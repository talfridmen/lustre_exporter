package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector interface {
	CollectBasicMetrics(ch chan<- prometheus.Metric)
	CollectExtendedMetrics(ch chan<- prometheus.Metric)
	Describe(ch chan<- *prometheus.Desc)
	GetLevel() Level
}

// Collector represents the interface for a collector
type BaseCollector struct {
	name  string
	level Level
}

func (c BaseCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	panic(fmt.Sprintf("CollectBasicMetrics is no implemented for collector %s", c.name))
}
func (c BaseCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	panic(fmt.Sprintf("CollectExtendedMetrics is no implemented for collector %s", c.name))
}
func (c BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	panic(fmt.Sprintf("Describe is no implemented for collector %s", c.name))
}

// GetLevel returns the level of operation for the collector
func (c *BaseCollector) GetLevel() Level {
	return c.level
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
		fmt.Printf("collector %s got unexpected level %s, disabling it.", collector, levelStr)
		return Disabled
	}
}
