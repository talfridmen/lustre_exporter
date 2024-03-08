package collectortypes

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/consts"
)

type SingleCollector struct {
	metric      *prometheus.Desc
	filePattern string
	level       consts.Level
}

func NewSingleCollector(metric *prometheus.Desc, filePattern string, level consts.Level) *SingleCollector {
	return &SingleCollector{
		metric:      metric,
		filePattern: filePattern,
		level:       level,
	}
}

func (x *SingleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.metric
}

func (c *SingleCollector) CollectSSingleMetric(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil || data == nil {
			fmt.Printf("could not read stat file %s\n", path)
		}
		value, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		ch <- prometheus.MustNewConstMetric(c.metric, prometheus.GaugeValue, float64(value), path)
	}
}

// CollectBasicMetrics collects basic metrics
func (c *SingleCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Basic {
		c.CollectSSingleMetric(ch, c.filePattern)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *SingleCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Extended {
		c.CollectSSingleMetric(ch, c.filePattern)
	}
}
