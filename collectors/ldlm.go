package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
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
					prometheus.NewDesc("lustre_ldlm_cancel_samples", "number of samples of metadata operations", []string{"path", "stat_type"}, nil),
					prometheus.NewDesc("lustre_ldlm_cancel_sum", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					prometheus.NewDesc("lustre_ldlm_cancel_sumsq", "number of samples of metadata operations", []string{"path", "stat_type", "units"}, nil),
					fmt.Sprintf("%s/ldlm/services/ldlm_canceld/stats", consts.KernelDebugBaseDir),
					consts.Extended,
				),
			},
			singleCollectors: []collectortypes.SingleCollector{
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_ldlm_lock_granted_count", "total number of locks granted", []string{"path"}, nil),
					fmt.Sprintf("%s/ldlm/lock_granted_count", consts.KernelDebugBaseDir),
					consts.Basic,
				),
				*collectortypes.NewSingleCollector(
					prometheus.NewDesc("lustre_ldlm_ost_lock_count", "total number of locks for an ost", []string{"path"}, nil),
					fmt.Sprintf("%s/ldlm/namespaces/filter-*/lock_count", consts.SysfsBaseDir),
					consts.Basic,
				),
			},
		},
	}
}
