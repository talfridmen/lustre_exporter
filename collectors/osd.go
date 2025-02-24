package collectors

import (
	"fmt"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	osdPathGlob = `osd-ldiskfs/*-*`
	osdPathReg  = `osd-ldiskfs/(?P<filesystem>.*)-(?P<osd>.*)`
)

type OsdCollector struct {
	BaseCollector
}

func NewOsdCollector(name string, level string) *OsdCollector {
	return &OsdCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_files_free", "number of free files in osd"),
					fmt.Sprintf("%s/%s/filesfree", consts.SysfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/filesfree`, consts.SysfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_files_total", "total number of files in osd"),
					fmt.Sprintf("%s/%s/filestotal", consts.SysfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/filestotal`, consts.SysfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_free", "free space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytesfree", consts.SysfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytesfree`, consts.SysfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_avail", "available space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytesavail", consts.SysfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytesavail`, consts.SysfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_total", "total space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytestotal", consts.SysfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytestotal`, consts.SysfsBaseDir, osdPathReg),
					consts.Basic,
				),
			},
			acctCollectors: []collectortypes.AcctCollector{
				*collectortypes.NewAcctCollector(
					collectortypes.NewMetricInfo("lustre_osd_acct_user_inodes", "inode accounting per user"),
					collectortypes.NewMetricInfo("lustre_osd_acct_user_kbytes", "size accounting per user in kbytes"),
					fmt.Sprintf("%s/%s/quota_slave/acct_user", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/quota_slave/acct_user`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewAcctCollector(
					collectortypes.NewMetricInfo("lustre_osd_acct_group_inodes", "inode accounting per group"),
					collectortypes.NewMetricInfo("lustre_osd_acct_group_kbytes", "size accounting per group in kbytes"),
					fmt.Sprintf("%s/%s/quota_slave/acct_group", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/quota_slave/acct_group`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewAcctCollector(
					collectortypes.NewMetricInfo("lustre_osd_acct_project_inodes", "inode accounting per project"),
					collectortypes.NewMetricInfo("lustre_osd_acct_project_kbytes", "size accounting per project in kbytes"),
					fmt.Sprintf("%s/%s/quota_slave/acct_project", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/quota_slave/acct_project`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
			},
		},
	}
}
