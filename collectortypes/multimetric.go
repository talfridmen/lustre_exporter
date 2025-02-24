package collectortypes

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	key   string
	value int
}

type MultiMetricCollector struct {
	metric      *prometheus.Desc
	filePattern string
	fileRegex   regexp.Regexp
	BaseCollector
}

func NewMultiMetricCollector(metric *MetricInfo, filePattern string, FileRegex string, configName string) *MultiMetricCollector {
	fileRegexp := *regexp.MustCompile(FileRegex)
	return &MultiMetricCollector{
		metric:      metric.CreatePrometheusMetric([]string{"key"}, fileRegexp),
		filePattern: filePattern,
		fileRegex:   fileRegexp,
		BaseCollector: BaseCollector{configKey: configName},
	}
}

func (x *MultiMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.metric
}

// CollectMetrics collects metrics
func (c *MultiMetricCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectMultiMetric(ch, c.filePattern)
}

func (c *MultiMetricCollector) CollectMultiMetric(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil || data == nil {
			fmt.Printf("could not read stat file %s\n", path)
		}

		metrics, err := ParseMetric(string(data))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, metric := range metrics {
			ch <- prometheus.MustNewConstMetric(c.metric, prometheus.GaugeValue, float64(metric.value), append([]string{metric.key}, c.fileRegex.FindStringSubmatch(path)[1:]...)...)
		}
	}
}

func ParseMetric(input string) ([]Metric, error) {
	var result []Metric

	scanner := bufio.NewScanner(strings.NewReader(input))

	metricLineRegex := regexp.MustCompile(`(?P<key>.*): (?P<value>\d*)`)

	for scanner.Scan() {
		line := scanner.Text()

		groups := metricLineRegex.FindStringSubmatch(line)

		value, err := strconv.Atoi(strings.TrimSpace(groups[2]))
		if err != nil {
			return nil, fmt.Errorf("error parsing value: %s (%s)", line, groups[2])
		}

		result = append(result, Metric{
			key:   groups[1],
			value: value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}
