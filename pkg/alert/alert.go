package alert

import (
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/collect/cmd"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/config"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/queue"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/util/parsers"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"time"
)

type Alert struct {
	MetricsQueue *queue.CollectMetricsQueue
	AlertQueue   *queue.AlertMetricsQueue
	Config       *config.AlertRuleConfig
	interval     time.Duration
}

func (c *Alert) diff(except config.AlertRule, current *queue.CollectMetrics) (*queue.AlertMetrics, bool) {
	expectVersion := ""
	result := &queue.AlertMetrics{
		Name:    except.Name,
		Current: current.Version,
		Cluster: c.Config.Cluster,
	}
	switch current.Type {
	case cmd.Tag:
		expectVersion = except.Version
		result.Type = cmd.Tag
	case cmd.Sha256:
		expectVersion = except.Sha265
		result.Type = cmd.Sha256
	}

	result.Expect = expectVersion
	return result, expectVersion == current.Version
}

func (c *Alert) run() {
	for _, except := range c.Config.List {
		if except.Name == "" {
			continue
		}
		for name, values := range c.MetricsQueue.List() {
			if name == except.Name {
				for _, current := range values {
					if r, ok := c.diff(except, current); !ok {
						klog.Warningf("app: %s, type: %s, current: %s, except: %s", name, r.Type, r.Current, r.Expect)
						c.AlertQueue.Set(except.Name, r)
					}
				}
			}
		}
	}
}

func (c *Alert) Run() {
	go wait.Forever(c.run, c.interval)
}

func NewAlert(cfg *config.AlertRuleConfig, metricsQueue *queue.CollectMetricsQueue, alertQueue *queue.AlertMetricsQueue) (*Alert, error) {
	interval, err := parsers.TimeParse(cfg.Interval)
	if err != nil {
		return nil, err
	}
	return &Alert{
		MetricsQueue: metricsQueue,
		AlertQueue:   alertQueue,
		Config:       cfg,
		interval:     interval,
	}, nil
}
