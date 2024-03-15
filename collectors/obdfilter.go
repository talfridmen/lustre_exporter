package collectors

import (
	"fmt"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	obdfilterPathGlob = `obdfilter/*`
	obdfilterPathReg  = `obdfilter/(?P<filesystem>.*)-(?P<ost>OST\d+)`
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
					collectortypes.NewMetricInfo("lustre_obdfilter_stats_sum", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_obdfilter_stats_sumsq", "number of samples of metadata operations"),
					fmt.Sprintf("%s/%s/stats", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/stats`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_export_stats_samples", "number of samples of metadata operations per export"),
					collectortypes.NewMetricInfo("lustre_obdfilter_export_stats_sum", "number of samples of metadata operations per export"),
					collectortypes.NewMetricInfo("lustre_obdfilter_export_stats_sumsq", "number of samples of metadata operations per export"),
					fmt.Sprintf("%s/%s/exports/*/stats", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/exports/(?P<ip>[\d\.]+)@(?P<network>.*)/stats`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_avail_kbytes", "ost space available for non-root users in kbytes"),
					fmt.Sprintf("%s/%s/kbytesavail", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/kbytesavail`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_free_kbytes", "ost free space in kbytes"),
					fmt.Sprintf("%s/%s/kbytesfree", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/kbytesfree`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_total_kbytes", "ost total space in kbytes"),
					fmt.Sprintf("%s/%s/kbytestotal", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/kbytestotal`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_files_free", "ost free inodes"),
					fmt.Sprintf("%s/%s/filesfree", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/filesfree`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_files_total", "ost total inodes"),
					fmt.Sprintf("%s/%s/filestotal", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/filestotal`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_brw_size_bytes", "ost bulk read/write size"),
					fmt.Sprintf("%s/%s/brw_size", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/brw_size`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_job_cleanup_interval_seconds", "the timeout for which an inactive job will be removed from memory"),
					fmt.Sprintf("%s/%s/job_cleanup_interval", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/job_cleanup_interval`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_exports_count", "number of times an ost is exported"),
					fmt.Sprintf("%s/%s/num_exports", consts.ProcfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/num_exports`, consts.ProcfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_obdfilter_degraded", "whether a pool is degraded or not"),
					fmt.Sprintf("%s/%s/degraded", consts.SysfsBaseDir, obdfilterPathGlob),
					fmt.Sprintf(`%s/%s/degraded`, consts.SysfsBaseDir, obdfilterPathReg),
					consts.Basic,
				),
			},
		},
	}
}
