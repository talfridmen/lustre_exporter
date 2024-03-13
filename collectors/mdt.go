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
					fmt.Sprintf("%s/mdt/*/md_stats", consts.ProcfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_mdt_export_stats_samples", "number of samples of metadata operations per export", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_mdt_export_stats_sum", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_mdt_export_stats_sumsq", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					fmt.Sprintf("%s/mdt/*/exports/*/stats", consts.ProcfsBaseDir),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_mdt_num_exports", "number f exports an mdt has", []string{"path"}, nil),
					fmt.Sprintf("%s/mdt/*/num_exports", consts.ProcfsBaseDir),
					consts.Basic,
				),
			},
			jobStatsCollectors: []collectortypes.JobStatsCollector{
				*collectortypes.NewJobStatsCollector(
					prometheus.NewDesc("lustre_mdt_job_stats_samples", "number of samples of metadata operations per job", []string{"path", "job", "stat_type"}, nil),
					prometheus.NewDesc("lustre_mdt_job_stats_sum", "number of samples of metadata operations per job", []string{"path", "job", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_mdt_job_stats_sumsq", "number of samples of metadata operations per job", []string{"path", "job", "stat_type", "units"}, nil),
					fmt.Sprintf("%s/mdt/*/job_stats", consts.ProcfsBaseDir),
					consts.Extended,
				),
			},
		},
	}
}
