package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

type MDTCollector struct {
	BaseCollector
}

func NewMDTCollector(name string, level string) *MDTCollector {
	return &MDTCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
			statsCollectors: []collectortypes.StatsCollector{
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_mdt_stats_samples", "number of samples of metadata operations", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_mdt_stats_sum", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_mdt_stats_sumsq", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					[]string{fmt.Sprintf("%s/mdt/*/md_stats", consts.ProcfsBaseDir)},
					[]string{},
				),
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_mdt_export_stats_samples", "number of samples of metadata operations per export", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_mdt_export_stats_sum", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_mdt_export_stats_sumsq", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					[]string{},
					[]string{fmt.Sprintf("%s/mdt/*/exports/*/stats", consts.ProcfsBaseDir)},
				),
			},
		},
	}
}
