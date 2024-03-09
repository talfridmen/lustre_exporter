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
					fmt.Sprintf("%s/obdfilter/*/stats", consts.ProcfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewStatsCollector(
					prometheus.NewDesc("lustre_obdfilter_export_stats_samples", "number of samples of metadata operations per export", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_obdfilter_export_stats_sum", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_obdfilter_export_stats_sumsq", "number of samples of metadata operations per export", []string{"path", "stat_type", "units"}, nil),
					fmt.Sprintf("%s/obdfilter/*/exports/*/stats", consts.ProcfsBaseDir),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_avail_kbytes", "ost space available for non-root users in kbytes", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/kbytesavail", consts.SysfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_free_kbytes", "ost free space in kbytes", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/kbytesfree", consts.SysfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_total_kbytes", "ost total space in kbytes", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/kbytestotal", consts.SysfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_files_free", "ost free inodes", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/filesfree", consts.SysfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_files_total", "ost total inodes", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/filestotal", consts.SysfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_brw_size_bytes", "ost bulk read/write size", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/brw_size", consts.ProcfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_job_cleanup_interval_seconds", "the timeout for which an inactive job will be removed from memory", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/job_cleanup_interval", consts.ProcfsBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_obdfilter_exports_count", "number of times an ost is exported", []string{"path"}, nil),
					fmt.Sprintf("%s/obdfilter/*/num_exports", consts.ProcfsBaseDir),
					consts.Basic,
				),
			},
		},
	}
}
