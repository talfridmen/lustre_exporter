package collectors

import (
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
}

func (c *BaseCollector) GetName() string {
	return c.name
}

func (c *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *BaseCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	for _, collector := range c.collectors {
		is_enabled, _ := c.config.Key(collector.GetConfigKey()).Bool()
		if is_enabled {
			collector.CollectMetrics(ch)
		}
	}
}
