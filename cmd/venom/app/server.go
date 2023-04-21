package app

import (
	"errors"
	"fmt"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/cmd/venom/app/options"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/alert"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/collect"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/config"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/prom"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/queue"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	componentName = "venom"
)

func NewVenomCommand() *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet(componentName, pflag.ContinueOnError)
	serverOption := options.NewVenomServer()
	cmd := &cobra.Command{
		Use:                componentName,
		Long:               `based on k8s to collect the version of each application`,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// initial flag parse, since we disable cobra's flag parsing
			if err := cleanFlagSet.Parse(args); err != nil {
				return fmt.Errorf("failed to parse %s flag: %w", componentName, err)
			}

			// check if there are non-flag arguments in the command line
			cmds := cleanFlagSet.Args()
			if len(cmds) > 0 {
				return fmt.Errorf("unknown command %+s", cmds[0])
			}

			// short-circuit on help
			help, err := cleanFlagSet.GetBool("help")
			if err != nil {
				return errors.New(`"help" flag is non-bool, programmer error, please correct`)
			}
			if help {
				return cmd.Help()
			}

			if serverOption.Cluster == "" {
				return errors.New("the region of a cluster must be filled in")
			}

			if serverOption.CollectRuleFile == "" {
				return errors.New("collection file does not exist")
			}

			cliflag.PrintFlags(cleanFlagSet)
			return run(serverOption)
		},
	}
	serverOption.AddFlags(cleanFlagSet)
	cleanFlagSet.BoolP("help", "h", false, fmt.Sprintf("help for %s", cmd.Name()))

	// ugly, but necessary, because Cobra's default UsageFunc and HelpFunc pollute the flagset with global flags
	const usageFmt = "Usage:\n  %s\n\nFlags:\n%s"
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine(), cleanFlagSet.FlagUsagesWrapped(2))
	})
	return cmd
}

func startHTTPServer(addr string, alertQ *queue.AlertMetricsQueue, metricsQ *queue.CollectMetricsQueue) {

	prom.RegisterPromCollector(alertQ, metricsQ)

	// http expose metrics
	router := gin.New()
	p := ginprometheus.NewPrometheus("venom")
	p.Use(router)
	prometheus.MustRegister(prom.NewCollector())

	router.GET("/healthy", func(ctx *gin.Context) {
		ctx.JSON(200, "OK")
	})

	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func run(server *options.VenomServer) error {
	alertRuleConfig, err := config.NewAlertRuleConfig(server.AlertRuleFile)
	if err != nil {
		return err
	}
	collectRuleConfig, err := config.NewCollectRuleConfig(server.CollectRuleFile)
	if err != nil {
		return err
	}

	metricsQueue, _ := queue.NewCollectMetricsQueue()
	alertQueue, _ := queue.NewAlertMetricsQueue()

	collectController, err := collect.NewCollect(collectRuleConfig, metricsQueue, server.CRIEngine)
	if err != nil {
		return err
	}

	alertController, err := alert.NewAlert(alertRuleConfig, metricsQueue, alertQueue)
	if err != nil {
		return nil
	}

	runf := func() {
		collectController.Run()
		alertController.Run()
		go startHTTPServer(server.WebListenAddress, alertQueue, metricsQueue)
	}

	// main service blocking signal
	stop := make(chan struct{})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		close(stop)
	}()

	runf()
	<-stop
	return nil
}
