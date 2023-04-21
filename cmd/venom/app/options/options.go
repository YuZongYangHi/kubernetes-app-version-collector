package options

import (
	"github.com/spf13/pflag"
)

type VenomFlags struct {
	Cluster          string
	CRIEngine        string
	WebListenAddress string
	AlertRuleFile    string
	CollectRuleFile  string
}

type VenomServer struct {
	VenomFlags
}

func (v *VenomFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()

	fs.StringVar(&v.WebListenAddress, "web-listen-address", "0.0.0.0:9001", "listening http address, usually used to expose metrics for prometheus")
	fs.StringVar(&v.Cluster, "cluster", "", "the name of the k8s cluster")
	fs.StringVar(&v.CRIEngine, "cri-engine", "docker", "container runtime")
	fs.StringVar(&v.AlertRuleFile, "alert-rule-file", "", "alert rule configuration file")
	fs.StringVar(&v.CollectRuleFile, "collect-rule-file", "", "collect rule configuration file")
}

func NewVenomServer() *VenomServer {
	return &VenomServer{
		VenomFlags{},
	}
}
