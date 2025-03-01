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

type RecoveryCollector struct {
	metric      *prometheus.Desc
	filePattern string
	fileRegex   regexp.Regexp
	BaseCollector
}

func NewRecoveryCollector(metric *MetricInfo, filePattern string, FileRegex string, configName string) *RecoveryCollector {
	fileRegexp := *regexp.MustCompile(FileRegex)
	return &RecoveryCollector{
		metric:      metric.CreatePrometheusMetric([]string{"key"}, fileRegexp),
		filePattern: filePattern,
		fileRegex:   fileRegexp,
		BaseCollector: BaseCollector{configKey: configName},
	}
}

func (x *RecoveryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.metric
}

// CollectMetrics collects metrics
func (c *RecoveryCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectRecoveryMetrics(ch, c.filePattern)
}

func (c *RecoveryCollector) CollectRecoveryMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil || data == nil {
			fmt.Printf("could not read recovery file %s\n", path)
		}

		metrics := ParseRecovery(string(data))
		for _, metric := range metrics {
			ch <- prometheus.MustNewConstMetric(c.metric, prometheus.GaugeValue, float64(metric.value), append([]string{metric.key}, c.fileRegex.FindStringSubmatch(path)[1:]...)...)
		}
	}
}

func ParseRecovery(input string) ([]Metric) {
	var result []Metric

	scanner := bufio.NewScanner(strings.NewReader(input))

	metricLineRegex := regexp.MustCompile(`(?P<key>.*): (?P<value>.*)`)

	for scanner.Scan() {
		line := scanner.Text()

		groups := metricLineRegex.FindStringSubmatch(line)
		key := strings.TrimSpace(groups[1])

		switch key {
		case "status":
			value := ParseStatus(groups[2])
			result = append(result, Metric{
				key:   key,
				value: value,
			})
		case "recovery_duration":
			value, _ := strconv.Atoi(strings.TrimSpace(groups[2]))
			result = append(result, Metric{
				key:   key,
				value: value,
			})
		case "completed_clients":
			completed, total := ParseClients(groups[2])
			result = append(
				result, 
				Metric{
					key:   fmt.Sprintf("%s_completed", key),
					value: completed,
				}, Metric{
					key:   fmt.Sprintf("%s_total", key),
					value: total,
				},
			)
		case "replayed_requests":
			value, _ := strconv.Atoi(strings.TrimSpace(groups[2]))
			result = append(result, Metric{
				key:   key,
				value: value,
			})
		}
	}

	return result
}

func ParseClients(value string) (completed, total int) {
	clientLineRegex := regexp.MustCompile(`(?P<completed>\d*)/(?P<total>\d*)`)
	groups := clientLineRegex.FindStringSubmatch(value)
	completed, _ = strconv.Atoi(strings.TrimSpace(groups[1]))
	total, _ = strconv.Atoi(strings.TrimSpace(groups[2]))

	return completed, total
}

var STATUSES = map[string]int{
	"COMPLETE": 0,
	"INACTIVE": 1,
	"WAITING": 2,
	"WAITING_FOR_CLIENTS": 3,
	"RECOVERING": 4,
}

func ParseStatus(value string) int {
	name := strings.TrimSpace(value)
	if val, ok := STATUSES[name]; ok {
		return val
	} else {
		return -1
	}
}