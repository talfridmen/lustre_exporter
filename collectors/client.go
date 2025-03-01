package collectors

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/talfridmen/lustre_exporter/collectortypes"
	"github.com/talfridmen/lustre_exporter/consts"
)

const (
	llitePathGlob = "llite/*"
	llitePathReg  = `llite/(?P<filesystem>.*)-[0-9a-fA-F]*`
)

type LliteCollector struct {
	BaseCollector
}

func NewLliteCollector(name string, config *ini.Section) *LliteCollector {
	return &LliteCollector{
		BaseCollector: BaseCollector{
			name:  name,
			config: *config,
			collectors: []collectortypes.CollectorType{
				collectortypes.NewMultiMetricCollector(
					collectortypes.NewMetricInfo("lustre_llite_cache", "info about lustre client cache"),
					fmt.Sprintf("%s/%s/max_cached_mb", consts.KernelDebugBaseDir, llitePathGlob),
					fmt.Sprintf("%s/%s/max_cached_mb", consts.KernelDebugBaseDir, llitePathReg),
					"cache",
				),
				collectortypes.NewStatsCollector(
					collectortypes.NewMetricInfo("lustre_llite_stats_samples", "number of samples of metadata operations"),
					collectortypes.NewMetricInfo("lustre_llite_stats_sum", "number of samples of metadata operations"),
					fmt.Sprintf("%s/%s/stats", consts.KernelDebugBaseDir, llitePathGlob),
					fmt.Sprintf(`%s/%s/stats`, consts.KernelDebugBaseDir, llitePathReg),
					"stats",
				),
			},
		},
	}
}
