package collectors

import (
	"fmt"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

type LdlmCollector struct {
	BaseCollector
}

func NewLdlmCollector(name string, level string) *LdlmCollector {
	return &LdlmCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
			statsCollectors: []collectortypes.StatsCollector{
				*collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_ldlm_cancel_samples", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_ldlm_cancel_sum", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_ldlm_cancel_sumsq", "number of samples of metadata operations"),
					fmt.Sprintf("%s/ldlm/services/ldlm_canceld/stats", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ldlm/services/ldlm_canceld/stats`, consts.KernelDebugBaseDir),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_ldlm_lock_granted_count", "total number of locks granted"),
					fmt.Sprintf("%s/ldlm/lock_granted_count", consts.KernelDebugBaseDir),
					fmt.Sprintf(`%s/ldlm/lock_granted_count`, consts.KernelDebugBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					collectortypes.NewMetricInfo("lustre_ldlm_ost_lock_count", "total number of locks for an ost"),
					fmt.Sprintf("%s/ldlm/namespaces/filter-*/lock_count", consts.SysfsBaseDir),
					fmt.Sprintf(`%s/ldlm/namespaces/filter-(?P<filesystem>.*)-(?P<ost>OST\d+)_UUID/lock_count`, consts.SysfsBaseDir),
					consts.Basic,
				),
			},
		},
	}
}
