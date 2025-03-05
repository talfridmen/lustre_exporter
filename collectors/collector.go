package collectors

import (
    "time"

	"gopkg.in/ini.v1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
)

type Collector interface {
	CollectMetrics(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
	GetName() string
}

// Collector represents the interface for a collector
type BaseCollector struct {
	name                  string
	config                ini.Section
	collectors            []collectortypes.CollectorType
	TimeMetric           *prometheus.Desc
}

func (c *BaseCollector) GetName() string {
	return c.name
}

func (c *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	c.TimeMetric = prometheus.NewDesc(
		"lustre_" + c.name + "_collection_time_seconds",
		"time in seconds to collect " + c.name + " metrics",
		nil,
		nil,
	)
	ch <- c.TimeMetric
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *BaseCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	startCollectionTime := time.Now()
	for _, collector := range c.collectors {
		is_enabled, _ := c.config.Key(collector.GetConfigKey()).Bool()
		if is_enabled {
			collector.CollectMetrics(ch)
		}
	}
	ch <- prometheus.MustNewConstMetric(c.TimeMetric, prometheus.GaugeValue, float64(time.Since(startCollectionTime)) / 1000000000)
}
