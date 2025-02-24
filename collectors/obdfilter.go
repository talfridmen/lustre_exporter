package collectors

import (
	"fmt"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	obdfilterPathGlob = `obdfilter/*`
	obdfilterPathReg  = `obdfilter/` + consts.OST_REG
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
					collectortypes.NewMetricInfo("lustre_obdfilter_stats_samples", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_obdfilter_stats_sum", "sum of sample sizes of metadata operations"),
					fmt.Sprintf("%s/%s/stats", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/stats`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_export_stats_samples", "number of samples of data operations per export"),
					collectortypes.NewMetricInfo("lustre_obdfilter_export_stats_sum", "sum of sample sizes of data operations per export"),
					fmt.Sprintf("%s/%s/exports/*/stats", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/exports/(?P<ip>[\d\.]+)@(?P<network>.*)/stats`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_brw_size_bytes", "ost bulk read/write size"),
					fmt.Sprintf("%s/%s/brw_size", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/brw_size`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_job_cleanup_interval_seconds", "the timeout for which an inactive job will be removed from memory"),
					fmt.Sprintf("%s/%s/job_cleanup_interval", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/job_cleanup_interval`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_exports_count", "number of times an ost is exported"),
					fmt.Sprintf("%s/%s/num_exports", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/num_exports`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_degraded", "whether a pool is degraded or not"),
					fmt.Sprintf("%s/%s/degraded", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/degraded`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
			},
			jobStatsCollectors: []collectortypes.JobStatsCollector{
				*collectortypes.NewJobStatsCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_job_stats_samples", "number of samples of data operations per job"),
					collectortypes.NewMetricInfo("lustre_obdfilter_job_stats_sum", "sum of sample sizes of data operations per job"),
					fmt.Sprintf("%s/%s/job_stats", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/job_stats`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Extended,
				),
			},
			quotaCollectors: []collectortypes.QuotaCollector{
				*collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_data_quota_hard_user", "hard quota per user"),
					collectortypes.NewMetricInfo("lustre_data_quota_soft_user", "soft quota per user"),
					fmt.Sprintf("%s/qmt/*/dt-0x0/glb-usr", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/dt-0x0/glb-usr`, consts.ProcfsBaseDir, consts.QMT_REG),
					consts.Basic,
				),
				*collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_data_quota_hard_group", "hard quota per group"),
					collectortypes.NewMetricInfo("lustre_data_quota_soft_group", "soft quota per group"),
					fmt.Sprintf("%s/qmt/*/dt-0x0/glb-grp", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/dt-0x0/glb-grp`, consts.ProcfsBaseDir, consts.QMT_REG),
					consts.Basic,
				),
				*collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_data_quota_hard_project", "hard quota per project"),
					collectortypes.NewMetricInfo("lustre_data_quota_soft_project", "soft quota per project"),
					fmt.Sprintf("%s/qmt/*/dt-0x0/glb-prj", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/dt-0x0/glb-prj`, consts.ProcfsBaseDir, consts.QMT_REG),
					consts.Basic,
				),
			},
		},
	}
}
