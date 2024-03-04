package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

type OBDFilterCollector struct {
	BaseCollector
}

func NewOBDFilterCollector(name string, level string) *OBDFilterCollector {
	return &OBDFilterCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
			statsCollectors: []collectortypes.StatsCollector{
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_obdfilter_stats_samples", "number of samples of metadata operations", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_obdfilter_stats_sum", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_obdfilter_stats_sumsq", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					[]string{fmt.Sprintf("%s/obdfilter/*/stats", consts.ProcfsBaseDir)},
					[]string{},
				),
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_obdfilter_export_stats_samples", "number of samples of metadata operations per export", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_obdfilter_export_stats_sum", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_obdfilter_export_stats_sumsq", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					[]string{},
					[]string{fmt.Sprintf("%s/obdfilter/*/exports/*/stats", consts.ProcfsBaseDir)},
				),
			},
		},
	}
}
