package collectors

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	mdtPathGlob = "mdt/*"
	mdtPathReg  = `mdt/` + consts.MDT_REG
)

type MDTCollector struct {
	BaseCollector
}

func NewMDTCollector(name string, config *ini.Section) *MDTCollector {
	return &MDTCollector{
		BaseCollector: BaseCollector{
			name:  name,
			config: *config,
			collectors: []collectortypes.CollectorType{
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_mdt_stats_samples", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_mdt_stats_sum", "sum of sample sizes of metadata operations"),
					fmt.Sprintf("%s/%s/md_stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/md_stats`, consts.ProcfsBaseDir, mdtPathReg),
					"stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_mdt_export_stats_samples", "number of samples of metadata operations per export"),
					collectortypes.NewMetricInfo("lustre_mdt_export_stats_sum", "sum of sample sizes of metadata operations per export"),
					fmt.Sprintf("%s/%s/exports/*/stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/exports/(?P<ip>[\d\.]+)@(?P<network>.*)/stats`, consts.ProcfsBaseDir, mdtPathReg),
					"exports",
				),
				collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_mdt_num_exports", "number of exports an mdt has"),
					fmt.Sprintf("%s/%s/num_exports", consts.SysfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/num_exports`, consts.SysfsBaseDir, mdtPathReg),
					"exports",
				),
				collectortypes.NewJobStatsCollector(
					collectortypes.NewMetricInfo("lustre_mdt_job_stats_samples", "number of samples of metadata operations per job"),
					collectortypes.NewMetricInfo("lustre_mdt_job_stats_sum", "sum of sample sizes of metadata operations per job"),
					fmt.Sprintf("%s/%s/job_stats", consts.ProcfsBaseDir, mdtPathGlob),
					fmt.Sprintf(`%s/%s/job_stats`, consts.ProcfsBaseDir, mdtPathReg),
					"jobstat",
				),
				collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_metadata_quota_hard_user", "hard quota per user"),
					collectortypes.NewMetricInfo("lustre_metadata_quota_soft_user", "soft quota per user"),
					fmt.Sprintf("%s/qmt/*/md-0x0/glb-usr", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/md-0x0/glb-usr`, consts.ProcfsBaseDir, consts.QMT_REG),
					"quota",
				),
				collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_metadata_quota_hard_group", "hard quota per group"),
					collectortypes.NewMetricInfo("lustre_metadata_quota_soft_group", "soft quota per group"),
					fmt.Sprintf("%s/qmt/*/md-0x0/glb-grp", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/md-0x0/glb-grp`, consts.ProcfsBaseDir, consts.QMT_REG),
					"quota",
				),
				collectortypes.NewQuotaCollector(
					collectortypes.NewMetricInfo("lustre_metadata_quota_hard_project", "hard quota per project"),
					collectortypes.NewMetricInfo("lustre_metadata_quota_soft_project", "soft quota per project"),
					fmt.Sprintf("%s/qmt/*/md-0x0/glb-prj", consts.ProcfsBaseDir),
					fmt.Sprintf(`%s/qmt/%s/md-0x0/glb-prj`, consts.ProcfsBaseDir, consts.QMT_REG),
					"quota",
				),
			},
		},
	}
}
