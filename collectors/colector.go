package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
)

type Collector interface {
	CollectBasicMetrics(ch chan<- prometheus.Metric)
	CollectExtendedMetrics(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
	GetLevel() Level
}

// Collector represents the interface for a collector
type BaseCollector struct {
	name            string
	level           Level
	statsCollectors []collectortypes.StatsCollector
}

// func (c BaseCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
// 	panic(fmt.Sprintf("CollectBasicMetrics is no implemented for collector %s", c.name))
// }
// func (c BaseCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
// 	panic(fmt.Sprintf("CollectExtendedMetrics is no implemented for collector %s", c.name))
// }
// func (c BaseCollector) Describe(ch chan<- *prometheus.Desc) {
// 	panic(fmt.Sprintf("Describe is no implemented for collector %s", c.name))
// }

// GetLevel returns the level of operation for the collector
func (c *BaseCollector) GetLevel() Level {
	return c.level
}

func (c *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *BaseCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectBasicMetrics(ch)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *BaseCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectExtendedMetrics(ch)
	}
}

// Level represents the operation level of a collector
type Level int

const (
	Disabled Level = iota
	Basic
	Extended
)

// getCollectorLevel determines the level for the collector based on user input
func getCollectorLevel(collector string, levelStr string) Level {
	switch levelStr {
	case "basic":
		return Basic
	case "extended":
		return Extended
	case "disabled":
		return Disabled
	default:
		fmt.Printf("collector %s got unexpected level %s, disabling it.\n", collector, levelStr)
		return Disabled
	}
}
