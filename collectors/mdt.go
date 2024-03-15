package collectors

import (
	"fmt"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	mdtPathGlob = "mdt/*"
	mdtPathReg  = `mdt/(?P<filsystem>.*)-(?P<mdt>MDT\d+)`
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
					collectortypes.NewMetricInfo("lustre_mdt_stats_samples", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_mdt_stats_sum", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_mdt_stats_sumsq", "number of samples of metadata operations"),
					fmt.Sprintf("%s/%s/md_stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/md_stats`, consts.ProcfsBaseDir, mdtPathReg),
					consts.Basic,
				),
				*collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_mdt_export_stats_samples", "number of samples of metadata operations per export"),
					collectortypes.NewMetricInfo("lustre_mdt_export_stats_sum", "number of samples of metadata operations per export"),
					collectortypes.NewMetricInfo("lustre_mdt_export_stats_sumsq", "number of samples of metadata operations per export"),
					fmt.Sprintf("%s/%s/exports/*/stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/exports/(?P<ip>[\d\.]+)@(?P<network>.*)/stats`, consts.ProcfsBaseDir, mdtPathReg),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_mdt_num_exports", "number f exports an mdt has"),
					fmt.Sprintf("%s/%s/num_exports", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/num_exports`, consts.ProcfsBaseDir, mdtPathReg),
					consts.Basic,
				),
			},
			jobStatsCollectors: []collectortypes.JobStatsCollector{
				*collectortypes.NewJobStatsCollector(
					collectortypes.NewMetricInfo("lustre_mdt_job_stats_samples", "number of samples of metadata operations per job"),
					collectortypes.NewMetricInfo("lustre_mdt_job_stats_sum", "number of samples of metadata operations per job"),
					collectortypes.NewMetricInfo("lustre_mdt_job_stats_sumsq", "number of samples of metadata operations per job"),
					fmt.Sprintf("%s/%s/job_stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/job_stats`, consts.ProcfsBaseDir, mdtPathReg),
					consts.Extended,
				),
			},
		},
	}
}
