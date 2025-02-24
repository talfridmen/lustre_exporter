package collectortypes

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type SingleCollector struct {
	metric      *prometheus.Desc
	filePattern string
	fileRegex   regexp.Regexp
	BaseCollector
}

func NewSingleCollector(metric *MetricInfo, filePattern string, FileRegex string, configName string) *SingleCollector {
	fileRegexp := *regexp.MustCompile(FileRegex)
	return &SingleCollector{
		metric:      metric.CreatePrometheusMetric([]string{}, fileRegexp),
		filePattern: filePattern,
		fileRegex:   fileRegexp,
		BaseCollector: BaseCollector{configKey: configName},
	}
}

func (x *SingleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.metric
}

// CollectMetrics collects metrics
func (c *SingleCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectSingleMetric(ch, c.filePattern)
}

func (c *SingleCollector) CollectSingleMetric(ch chan<- prometheus.Metric, pattern string) {
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
		ch <- prometheus.MustNewConstMetric(c.metric, prometheus.GaugeValue, float64(value), c.fileRegex.FindStringSubmatch(path)[1:]...)
	}
}