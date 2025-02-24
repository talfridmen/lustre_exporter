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
	singleCollectors      []collectortypes.SingleCollector
	multiMetricCollectors []collectortypes.MultiMetricCollector
	statsCollectors       []collectortypes.StatsCollector
	jobStatsCollectors    []collectortypes.JobStatsCollector
	quotaCollectors       []collectortypes.QuotaCollector
	acctCollectors        []collectortypes.AcctCollector
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
func (c *BaseCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	for _, singleCollector := range c.singleCollectors {
		is_enabled, _ := c.config.Key(singleCollector.GetConfigKey()).Bool()
		if is_enabled {
			singleCollector.CollectMetrics(ch)
		}
	}
	for _, multiMetricCollector := range c.multiMetricCollectors {
		is_enabled, _ := c.config.Key(multiMetricCollector.GetConfigKey()).Bool()
		if is_enabled {
			multiMetricCollector.CollectMetrics(ch)
		}
		
	}
	for _, statsCollector := range c.statsCollectors {
		is_enabled, _ := c.config.Key(statsCollector.GetConfigKey()).Bool()
		if is_enabled {
			statsCollector.CollectMetrics(ch)
		}
	}
	for _, jobStatsCollector := range c.jobStatsCollectors {
		is_enabled, _ := c.config.Key(jobStatsCollector.GetConfigKey()).Bool()
		if is_enabled {
			jobStatsCollector.CollectMetrics(ch)
		}
	}
	for _, quotaCollector := range c.quotaCollectors {
		is_enabled, _ := c.config.Key(quotaCollector.GetConfigKey()).Bool()
		if is_enabled {
			quotaCollector.CollectMetrics(ch)
		}
	}
	for _, acctCollector := range c.acctCollectors {
		is_enabled, _ := c.config.Key(acctCollector.GetConfigKey()).Bool()
		if is_enabled {
			acctCollector.CollectMetrics(ch)
		}
	}
}
