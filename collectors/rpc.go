package collectors

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

type RpcCollector struct {
	BaseCollector
}

func NewRpcCollector(name string, config *ini.Section) *RpcCollector {
	return &RpcCollector{
		BaseCollector: BaseCollector{
			name:  name,
			config: *config,
			collectors: []collectortypes.CollectorType{
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_stats_samples", "number of samples of metadata rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_stats_sum", "sum of metadata rpcs"),
					fmt.Sprintf("%s/mds/MDS/mdt/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/mds/MDS/mdt/stats`, consts.KernelDebugBaseDir),
					"mds_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_io_stats_samples", "number of samples of metadata rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_io_stats_sum", "sum of metadata rpcs"),
					fmt.Sprintf("%s/mds/MDS/mdt_io/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/mds/MDS/mdt_io/stats`, consts.KernelDebugBaseDir),
					"mds_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_out_stats_samples", "number of samples of metadata rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_out_stats_sum", "sum of metadata rpcs"),
					fmt.Sprintf("%s/mds/MDS/mdt_out/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/mds/MDS/mdt_out/stats`, consts.KernelDebugBaseDir),
					"mds_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_fld_stats_samples", "number of samples of metadata rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_mds_mdt_fld_stats_sum", "sum of metadata rpcs"),
					fmt.Sprintf("%s/mds/MDS/mdt_fld/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/mds/MDS/mdt_fld/stats`, consts.KernelDebugBaseDir),
					"mds_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_stats_samples", "number of samples of data rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_stats_sum", "sum of data rpcs"),
					fmt.Sprintf("%s/ost/OSS/ost/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ost/OSS/ost/stats`, consts.KernelDebugBaseDir),
					"oss_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_io_stats_samples", "number of samples of data rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_io_stats_sum", "sum of data rpcs"),
					fmt.Sprintf("%s/ost/OSS/ost_io/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ost/OSS/ost_io/stats`, consts.KernelDebugBaseDir),
					"oss_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_out_stats_samples", "number of samples of data rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_out_stats_sum", "sum of data rpcs"),
					fmt.Sprintf("%s/ost/OSS/ost_out/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ost/OSS/ost_out/stats`, consts.KernelDebugBaseDir),
					"oss_stats",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_create_stats_samples", "number of samples of data rpcs"),
					collectortypes.NewMetricInfo("lustre_rpc_oss_ost_create_stats_sum", "sum of data rpcs"),
					fmt.Sprintf("%s/ost/OSS/ost_create/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ost/OSS/ost_create/stats`, consts.KernelDebugBaseDir),
					"oss_stats",
				),
				collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_mds_threads", "total number of threads active for metadata operations"),
					fmt.Sprintf("%s/mds/MDS/mdt/threads_started", consts.SysfsBaseDir),
					fmt.Sprintf(`%s/mds/MDS/mdt/threads_started`, consts.SysfsBaseDir),
					"threads",
				),
				collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_oss_threads", "total number of threads active for data operations"),
					fmt.Sprintf("%s/ost/OSS/ost/threads_started", consts.SysfsBaseDir),
					fmt.Sprintf(`%s/ost/OSS/ost/threads_started`, consts.SysfsBaseDir),
					"threads",
				),
			},
		},
	}
}
