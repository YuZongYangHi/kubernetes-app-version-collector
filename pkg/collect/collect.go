package collect

import (
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/collect/cmd"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/config"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/queue"
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/util/parsers"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"time"
)

type Collect struct {
	Config   *config.CollectRuleConfig
	Cmd      *cmd.Cmd
	Queue    *queue.CollectMetricsQueue
	interval time.Duration
}

func (c *Collect) Run() {
	go wait.Forever(c.run, c.interval)
}

func (c *Collect) run() {
	for _, item := range c.Config.List {
		if item.Name == "" {
			continue
		}

		result, err := c.Cmd.Run(item.Name, item.Cmd)
		if err != nil {
			klog.Errorf("collect app: %s fail: %s", item.Name, err.Error())
			continue
		}
		c.Queue.Set(item.Name, result)
	}
}

func NewCollect(cfg *config.CollectRuleConfig, queue *queue.CollectMetricsQueue, cri string) (*Collect, error) {
	c, err := cmd.NewCmd(cri)
	if err != nil {
		return nil, err
	}
	interval, err := parsers.TimeParse(cfg.Interval)
	if err != nil {
		return nil, err
	}

	return &Collect{
		Config:   cfg,
		Cmd:      c,
		Queue:    queue,
		interval: interval,
	}, nil
}
