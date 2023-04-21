package prom

import (
	"github.com/YuZongYangHi/kubernetes-app-version-collector/pkg/queue"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	collectMetricsMetadata = []string{"name", "type", "version", "cluster"}
	alertMetricsMetadata   = []string{"name", "type", "current", "expect", "cluster"}
)

func RegisterPromCollector(alertQ *queue.AlertMetricsQueue, metricsQ *queue.CollectMetricsQueue) {
	RegisterCollector(NewVersionCollectMetrics(metricsQ), "venom_version_collect")
	RegisterCollector(NewVersionCollectDiffMetrics(alertQ), "venom_version_diff")
}

type VersionCollectDiffMetrics struct {
	alertQ *queue.AlertMetricsQueue
	name   *prometheus.Desc
}

type VersionCollectMetrics struct {
	metricsQ *queue.CollectMetricsQueue
	name     *prometheus.Desc
}

func NewVersionCollectMetrics(q *queue.CollectMetricsQueue) *VersionCollectMetrics {
	return &VersionCollectMetrics{
		metricsQ: q,
		name: prometheus.NewDesc("venom_collector_version",
			"venom service app collector version",
			collectMetricsMetadata, nil,
		),
	}
}

func NewVersionCollectDiffMetrics(q *queue.AlertMetricsQueue) *VersionCollectDiffMetrics {
	return &VersionCollectDiffMetrics{
		alertQ: q,
		name: prometheus.NewDesc("venom_version_diff",
			"venom service collect version diff",
			alertMetricsMetadata, nil,
		),
	}
}

func (t *VersionCollectMetrics) Update(c chan<- prometheus.Metric) error {
	for _, values := range t.metricsQ.List() {
		for _, value := range values {
			c <- prometheus.MustNewConstMetric(t.name, prometheus.GaugeValue, 0,
				value.Name, value.Type, value.Version, value.Cluster,
			)
		}
	}
	return nil
}

func (t *VersionCollectDiffMetrics) Update(c chan<- prometheus.Metric) error {
	for _, value := range t.alertQ.List() {
		c <- prometheus.MustNewConstMetric(t.name, prometheus.GaugeValue, 1,
			value.Name, value.Type, value.Current, value.Expect, value.Cluster,
		)
	}

	return nil
}

func (t *VersionCollectMetrics) MetricName() string {
	return ""
}

func (t *VersionCollectDiffMetrics) MetricName() string {
	return ""
}
