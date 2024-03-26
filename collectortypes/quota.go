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
	"github.com/talfridmen/lustre_exporter/consts"
)

type QuotaCollector struct {
	QuotaHardMetric   *prometheus.Desc
	QuotaSoftMetric   *prometheus.Desc
	quotaFilePatterns string
	quotaFileRegex    regexp.Regexp
	level             consts.Level
}

// SampleData represents the parsed information for each label
type QuotaEntry struct {
	Id        string
	hardLimit int
	softLimit int
	granted   string
	time      string
}

func NewQuotaCollector(QuotaHardMetric *MetricInfo, QuotaSoftMetric *MetricInfo, quotaFilePatterns string, quotaFileRegex string, level consts.Level) *QuotaCollector {
	quotaFileRegexp := *regexp.MustCompile(quotaFileRegex)
	return &QuotaCollector{
		QuotaHardMetric:   QuotaHardMetric.CreatePrometheusMetric([]string{"id", "granted", "time"}, quotaFileRegexp),
		QuotaSoftMetric:   QuotaSoftMetric.CreatePrometheusMetric([]string{"id", "granted", "time"}, quotaFileRegexp),
		quotaFilePatterns: quotaFilePatterns,
		quotaFileRegex:    quotaFileRegexp,
		level:             level,
	}
}

func (x *QuotaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.QuotaHardMetric
	ch <- x.QuotaSoftMetric
}

// CollectBasicMetrics collects basic metrics
func (c *QuotaCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Basic {
		c.CollectQuotaMetrics(ch, c.quotaFilePatterns)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *QuotaCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Extended {
		c.CollectQuotaMetrics(ch, c.quotaFilePatterns)
	}
}

func (c *QuotaCollector) CollectQuotaMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		pathLabels := c.quotaFileRegex.FindStringSubmatch(path)[1:]
		value, err := os.ReadFile(filepath.Clean(path))
		if err != nil || value == nil {
			fmt.Printf("could not read quota file %s\n", path)
		}
		quotas, err := ParseQuotas(string(value))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, quota := range quotas {
			ch <- prometheus.MustNewConstMetric(c.QuotaHardMetric, prometheus.GaugeValue, float64(quota.hardLimit), append([]string{quota.Id, quota.granted, quota.time}, pathLabels...)...)
			ch <- prometheus.MustNewConstMetric(c.QuotaSoftMetric, prometheus.GaugeValue, float64(quota.softLimit), append([]string{quota.Id, quota.granted, quota.time}, pathLabels...)...)
		}
	}
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseQuotas(input string) (map[string]QuotaEntry, error) {
	var result map[string]QuotaEntry = make(map[string]QuotaEntry)

	scanner := bufio.NewScanner(strings.NewReader(input))
	id := ""
	scanner.Scan()
	scanner.Text()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "- id:") {
			suffix, _ := strings.CutPrefix(line, "- id:")
			id = strings.TrimSpace(suffix)
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 11 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		hardLimit, err := strconv.Atoi(strings.Trim(fields[3], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse hard limit: %v", err)
		}

		softLimit, err := strconv.Atoi(strings.Trim(fields[5], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse soft limit: %v", err)
		}

		granted := strings.Trim(fields[7], ",")

		time := strings.Trim(fields[9], ",")

		result[id] = QuotaEntry{
			Id:        id,
			hardLimit: hardLimit,
			softLimit: softLimit,
			granted:   granted,
			time:      time,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}
