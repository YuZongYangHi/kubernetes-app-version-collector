package queue

import (
	"fmt"
	"sync"
)

type CollectMetricsQueue struct {
	data sync.Map
}

type AlertMetricsQueue struct {
	data sync.Map
}

func (c *CollectMetricsQueue) Set(key string, value []*CollectMetrics) {
	c.data.Store(key, value)
}

func (c *CollectMetricsQueue) List() map[string][]*CollectMetrics {
	result := make(map[string][]*CollectMetrics)

	c.data.Range(func(key, value any) bool {
		result[key.(string)] = value.([]*CollectMetrics)
		return true
	})
	return result
}

func (c *AlertMetricsQueue) key(key string, value *AlertMetrics) string {
	return fmt.Sprintf("%s_%s", key, value.Type)
}

func (c *AlertMetricsQueue) Set(key string, value *AlertMetrics) {
	c.data.Store(c.key(key, value), value)
}

func (c *AlertMetricsQueue) Delete(key string, value *AlertMetrics) {
	c.data.Delete(c.key(key, value))
}

func (c *AlertMetricsQueue) List() []*AlertMetrics {
	var result []*AlertMetrics
	c.data.Range(func(key, value any) bool {
		result = append(result, value.(*AlertMetrics))
		return true
	})
	return result
}

func NewCollectMetricsQueue() (*CollectMetricsQueue, error) {
	return &CollectMetricsQueue{}, nil
}

func NewAlertMetricsQueue() (*AlertMetricsQueue, error) {
	return &AlertMetricsQueue{}, nil
}
