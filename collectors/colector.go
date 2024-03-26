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
	GetName() string
}

// Collector represents the interface for a collector
type BaseCollector struct {
	name                  string
	level                 consts.Level
	singleCollectors      []collectortypes.SingleCollector
	multiMetricCollectors []collectortypes.MultiMetricCollector
	statsCollectors       []collectortypes.StatsCollector
	jobStatsCollectors    []collectortypes.JobStatsCollector
	quotaCollectors       []collectortypes.QuotaCollector
	acctCollectors        []collectortypes.AcctCollector
}

// GetLevel returns the level of operation for the collector
func (c *BaseCollector) GetLevel() consts.Level {
	return c.level
}

func (c *BaseCollector) GetName() string {
	return c.name
}

func (c *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, singleCollector := range c.singleCollectors {
		singleCollector.Describe(ch)
	}
	for _, multiMetricCollector := range c.multiMetricCollectors {
		multiMetricCollector.Describe(ch)
	}
	for _, statsCollector := range c.statsCollectors {
		statsCollector.Describe(ch)
	}
	for _, jobStatsCollector := range c.jobStatsCollectors {
		jobStatsCollector.Describe(ch)
	}
	for _, quotaCollector := range c.quotaCollectors {
		quotaCollector.Describe(ch)
	}
	for _, acctCollector := range c.acctCollectors {
		acctCollector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *BaseCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	for _, singleCollector := range c.singleCollectors {
		singleCollector.CollectBasicMetrics(ch)
	}
	for _, multiMetricCollector := range c.multiMetricCollectors {
		multiMetricCollector.CollectBasicMetrics(ch)
	}
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectBasicMetrics(ch)
	}
	for _, jobStatsCollector := range c.jobStatsCollectors {
		jobStatsCollector.CollectBasicMetrics(ch)
	}
	for _, quotaCollector := range c.quotaCollectors {
		quotaCollector.CollectBasicMetrics(ch)
	}
	for _, acctCollector := range c.acctCollectors {
		acctCollector.CollectBasicMetrics(ch)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *BaseCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	for _, singleCollector := range c.singleCollectors {
		singleCollector.CollectExtendedMetrics(ch)
	}
	for _, multiMetricCollector := range c.multiMetricCollectors {
		multiMetricCollector.CollectExtendedMetrics(ch)
	}
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectExtendedMetrics(ch)
	}
	for _, jobStatsCollector := range c.jobStatsCollectors {
		jobStatsCollector.CollectExtendedMetrics(ch)
	}
	for _, quotaCollector := range c.quotaCollectors {
		quotaCollector.CollectExtendedMetrics(ch)
	}
	for _, acctCollector := range c.acctCollectors {
		acctCollector.CollectExtendedMetrics(ch)
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
