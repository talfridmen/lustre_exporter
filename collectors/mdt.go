package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	mdt_basic_stats_file_patterns    = [...]string{"mdt/*/md_stats"}
	mdt_extended_stats_file_patterns = [...]string{}
)

type MDTCollector struct {
	BaseCollector
	statsCollector
}

func NewMDTCollector(name string, level string) *MDTCollector {
	return &MDTCollector{
		BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
		},
		statsCollector{
			stats_samples_metric: prometheus.NewDesc(
				"lustre_mdt_stats_samples",
				"number of samples of metadata operations",
				[]string{"mdt", "stat_type"},
				nil,
			),
			stats_sum_metric: prometheus.NewDesc(
				"lustre_mdt_stats_sum",
				"number of samples of metadata operations",
				[]string{"mdt", "stat_type", "units"},
				nil,
			),
			stats_sumsq_metric: prometheus.NewDesc(
				"lustre_mdt_stats_sumsq",
				"number of samples of metadata operations",
				[]string{"mdt", "stat_type", "units"},
				nil,
			),
			basic_stats_file_patterns:    mdt_basic_stats_file_patterns[:],
			extended_stats_file_patterns: mdt_extended_stats_file_patterns[:],
		},
	}
}

func (x *MDTCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.stats_samples_metric
	ch <- x.stats_sum_metric
	ch <- x.stats_sumsq_metric
}

// CollectBasicMetrics collects basic metrics
func (c *MDTCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	c.statsCollector.CollectBasicMetrics(ch)
}

// CollectExtendedMetrics collects extended metrics
func (c *MDTCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	c.statsCollector.CollectExtendedMetrics(ch)
}
