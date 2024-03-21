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
					fmt.Sprintf("%s/%s/filesfree", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/filesfree`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_files_total", "total number of files in osd"),
					fmt.Sprintf("%s/%s/filestotal", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/filestotal`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_free", "free space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytesfree", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytesfree`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_avail", "available space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytesavail", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytesavail`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_osd_kbytes_total", "total space in osd in kbytes"),
					fmt.Sprintf("%s/%s/kbytestotal", consts.ProcfsBaseDir, osdPathGlob),
					fmt.Sprintf(`%s/%s/kbytestotal`, consts.ProcfsBaseDir, osdPathReg),
					consts.Basic,
				),
			},
		},
	}
}
