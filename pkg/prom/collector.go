package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

var (
	collectorApps      = make(map[string]MetricsFactory)
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName("venom", "scrape", "collector_duration_seconds"),
		"venom: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName("venom", "scrape", "collector_success"),
		"venom: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

type MetricsFactory interface {
	Update(ch chan<- prometheus.Metric) error
	MetricName() string
}

type Collector struct {
	Collectors map[string]MetricsFactory
}

func RegisterCollector(app MetricsFactory, name string) {
	collectorApps[name] = app
}

func NewCollector() *Collector {
	c := new(Collector)
	c.Collectors = make(map[string]MetricsFactory)
	for name, app := range collectorApps {
		c.Collectors[name] = app
	}
	return c
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	wg := sync.WaitGroup{}
	wg.Add(len(collector.Collectors))
	for name, c := range collector.Collectors {
		go func(name string, c MetricsFactory) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c MetricsFactory, ch chan<- prometheus.Metric) {
	begin := time.Now()
	duration := time.Since(begin)
	var success float64

	if err := c.Update(ch); err != nil {
		success = 0
	} else {
		success = 1
	}

	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}
