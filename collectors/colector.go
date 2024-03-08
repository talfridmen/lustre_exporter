package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

type Collector interface {
	CollectBasicMetrics(ch chan<- prometheus.Metric)
	CollectExtendedMetrics(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
	GetLevel() consts.Level
}

// Collector represents the interface for a collector
type BaseCollector struct {
	name             string
	level            consts.Level
	statsCollectors  []collectortypes.StatsCollector
	singleCollectors []collectortypes.SingleCollector
}

// GetLevel returns the level of operation for the collector
func (c *BaseCollector) GetLevel() consts.Level {
	return c.level
}

func (c *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.Describe(ch)
	}
	for _, singleCollector := range c.singleCollectors {
		singleCollector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *BaseCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectBasicMetrics(ch)
	}
	for _, singleCollector := range c.singleCollectors {
		singleCollector.CollectBasicMetrics(ch)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *BaseCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectExtendedMetrics(ch)
	}
	for _, singleCollector := range c.singleCollectors {
		singleCollector.CollectExtendedMetrics(ch)
	}
}

// getCollectorLevel determines the level for the collector based on user input
func getCollectorLevel(collector string, levelStr string) consts.Level {
	switch levelStr {
	case "basic":
		return consts.Basic
	case "extended":
		return consts.Extended
	case "disabled":
		return consts.Disabled
	default:
		fmt.Printf("collector %s got unexpected level %s, disabling it.\n", collector, levelStr)
		return consts.Disabled
	}
}
