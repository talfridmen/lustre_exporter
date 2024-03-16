package collectors

import (
	"fmt"

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

func NewLliteCollector(name string, level string) *LliteCollector {
	return &LliteCollector{
		BaseCollector: BaseCollector{
			name:  name,
			level: getCollectorLevel(name, level),
			multiMetricCollectors: []collectortypes.MultiMetricCollector{
				*collectortypes.NewMultiMetricCollector(
					collectortypes.NewMetricInfo("lustre_llite_cache", "info about lustre client cache"),
					fmt.Sprintf("%s/%s/max_cached_mb", consts.KernelDebugBaseDir, llitePathGlob),
					fmt.Sprintf("%s/%s/max_cached_mb", consts.KernelDebugBaseDir, llitePathReg),
					consts.Basic,
				),
			},
		},
	}
}
