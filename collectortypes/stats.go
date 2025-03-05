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

// SampleData represents the parsed information for each label
type Stat struct {
	Syscall    string
	NumSamples int
	Unit       string
	Sum        int
}

type StatsCollector struct {
	statsSamplesMetric *prometheus.Desc
	statsSumMetric     *prometheus.Desc
	statsFilePatterns  string
	statsFileRegex     regexp.Regexp
	BaseCollector
}

func NewStatsCollector(statsSamplesMetric *MetricInfo, statsSumMetric *MetricInfo, statsFilePatterns string, statsFileRegex string, configName string) *StatsCollector {
	statsFileRegexp := *regexp.MustCompile(statsFileRegex)
	return &StatsCollector{
		statsSamplesMetric: statsSamplesMetric.CreatePrometheusMetric([]string{"syscall"}, statsFileRegexp),
		statsSumMetric:     statsSumMetric.CreatePrometheusMetric([]string{"syscall", "units"}, statsFileRegexp),
		statsFilePatterns:  statsFilePatterns,
		statsFileRegex:     statsFileRegexp,
		BaseCollector: BaseCollector{configKey: configName},
	}
}

func (x *StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.statsSamplesMetric
	ch <- x.statsSumMetric
}

// CollectMetrics collects metrics
func (c *StatsCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.statsFilePatterns)
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseStats(input string) (map[string]Stat, error) {
	var result map[string]Stat = make(map[string]Stat)

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 4 {
			// Skip some time related lines
			continue
		}

		syscall := fields[0]

		numSamples, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of samples: %v", err)
		}

		unit := fields[3][1 : len(fields[3])-1] // Extracting unit from [usecs]

		if len(fields) < 7 {
			result[syscall] = Stat{
				Syscall:    syscall,
				NumSamples: numSamples,
				Unit:       unit,
			}
			continue
		}

		sum, err := strconv.Atoi(fields[6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum value: %v", err)
		}

		result[syscall] = Stat{
			Syscall:    syscall,
			NumSamples: numSamples,
			Unit:       unit,
			Sum:        sum,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}


func (c *StatsCollector) CollectStatMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		pathLabels := c.statsFileRegex.FindStringSubmatch(path)[1:]
		value, err := os.ReadFile(filepath.Clean(path))
		if err != nil || value == nil {
			fmt.Printf("could not read stat file %s\n", path)
		}
		stats, err := ParseStats(string(value))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, stat := range stats {
			ch <- prometheus.MustNewConstMetric(c.statsSamplesMetric, prometheus.GaugeValue, float64(stat.NumSamples), append([]string{stat.Syscall}, pathLabels...)...)
			ch <- prometheus.MustNewConstMetric(c.statsSumMetric, prometheus.GaugeValue, float64(stat.Sum), append([]string{stat.Syscall, stat.Unit}, pathLabels...)...)
		}
	}
}
