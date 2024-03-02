package collectortypes

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/talfridmen/lustre_exporter/consts"
)

// TODO: remove after use
// var (
// 	basic_stats_file_patterns    = [...]string{"mdt/*/md_stats", "obdfilter/*/stats"}
// 	extended_stats_file_patterns = [...]string{"mdt/*/exports/*/stats", "obdfilter/*/exports/*/stats", "ldlm.namespaces.filter-*.pool.stats"}
// )

// SampleData represents the parsed information for each label
type Stat struct {
	Syscall                   string
	NumSamples                int
	Unit                      string
	Min, Max, Sum, SumSquared int
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseStats(input string) ([]Stat, error) {
	var result []Stat

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "snapshot_time") {
			// Skip the snapshot_time line
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		syscall := fields[0]

		numSamples, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of samples: %v", err)
		}

		unit := fields[3][1 : len(fields[3])-1] // Extracting unit from [usecs]
		min, err := strconv.Atoi(fields[4])
		if err != nil {
			return nil, fmt.Errorf("failed to parse min value: %v", err)
		}

		max, err := strconv.Atoi(fields[5])
		if err != nil {
			return nil, fmt.Errorf("failed to parse max value: %v", err)
		}

		sum, err := strconv.Atoi(fields[6])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum value: %v", err)
		}

		sumSquared, err := strconv.Atoi(fields[7])
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum squared value: %v", err)
		}

		result = append(result, Stat{
			Syscall:    syscall,
			NumSamples: numSamples,
			Unit:       unit,
			Min:        min,
			Max:        max,
			Sum:        sum,
			SumSquared: sumSquared,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}

type StatsCollector struct {
	statsSamplesMetric        *prometheus.Desc
	statsSumMetric            *prometheus.Desc
	statsSumsqMetric          *prometheus.Desc
	basicStatsFilePatterns    []string
	extendedStatsFilePatterns []string
}

func NewStatsCollector(statsSamplesMetric *prometheus.Desc, statsSumMetric *prometheus.Desc, statsSumsqMetric *prometheus.Desc, basicStatsFilePatterns []string, extendedStatsFilePatterns []string) *StatsCollector {
	return &StatsCollector{
		statsSamplesMetric:        statsSamplesMetric,
		statsSumMetric:            statsSumMetric,
		statsSumsqMetric:          statsSumsqMetric,
		basicStatsFilePatterns:    basicStatsFilePatterns,
		extendedStatsFilePatterns: extendedStatsFilePatterns,
	}
}

func (x *StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.statsSamplesMetric
	ch <- x.statsSumMetric
	ch <- x.statsSumsqMetric
}

func (c *StatsCollector) CollectStatMetrics(ch chan<- prometheus.Metric, patterns []string) {
	for _, pattern := range patterns {
		fullPath := filepath.Join(consts.ProcfsBaseDir, pattern)
		paths, _ := filepath.Glob(fullPath)
		if paths == nil {
			continue
		}
		for _, path := range paths {
			patternPath := strings.Replace(path, consts.ProcfsBaseDir, "", 1)
			value, err := os.ReadFile(filepath.Clean(path))
			if err != nil || value == nil {
				fmt.Printf("could not read stat file %s\n", path)
			}
			stats, err := ParseStats(string(value))
			if err != nil {
				fmt.Printf("got error while parsing line: %s\n", err)
			}
			for _, stat := range stats {
				ch <- prometheus.MustNewConstMetric(c.statsSamplesMetric, prometheus.GaugeValue, float64(stat.NumSamples), patternPath, stat.Syscall)
				ch <- prometheus.MustNewConstMetric(c.statsSumMetric, prometheus.GaugeValue, float64(stat.Sum), patternPath, stat.Syscall, stat.Unit)
				ch <- prometheus.MustNewConstMetric(c.statsSumsqMetric, prometheus.GaugeValue, float64(stat.SumSquared), patternPath, stat.Syscall, stat.Unit)
			}
		}
	}
}

// CollectBasicMetrics collects basic metrics
func (c *StatsCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.basicStatsFilePatterns[:])
}

// CollectExtendedMetrics collects extended metrics
func (c *StatsCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.extendedStatsFilePatterns[:])
}
