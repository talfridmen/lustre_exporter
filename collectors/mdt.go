package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
)

var (
	mdtBasicStatsFilePatterns    = [...]string{"mdt/*/md_stats"}
	mdtExtendedStatsFilePatterns = [...]string{}
)

type MDTCollector struct {
	BaseCollector
	statsCollectors []collectortypes.StatsCollector
}

func NewMDTCollector(name string, level string) *MDTCollector {
	return &MDTCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
		},

		statsCollectors: []collectortypes.StatsCollector{
			*collectortypes.NewStatsCollector(
				prometheus.NewDesc("lustre_mdt_stats_samples", "number of samples of metadata operations", []string{"path", "stat_type"}, nil),
				prometheus.NewDesc("lustre_mdt_stats_sum", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
				prometheus.NewDesc("lustre_mdt_stats_sumsq", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
				mdtBasicStatsFilePatterns[:],
				mdtExtendedStatsFilePatterns[:]),
		},
	}
}

func (c *MDTCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.Describe(ch)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *MDTCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectBasicMetrics(ch)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *MDTCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	for _, statsCollector := range c.statsCollectors {
		statsCollector.CollectExtendedMetrics(ch)
	}
}
