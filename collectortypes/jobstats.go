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
type JobStat struct {
	Job                       string
	Syscall                   string
	NumSamples                int
	Unit                      string
	Min, Max, Sum, SumSquared int
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseJobStat(input string) ([]JobStat, error) {
	var result []JobStat

	scanner := bufio.NewScanner(strings.NewReader(input))
	job := ""
	scanner.Scan()
	scanner.Text()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "snapshot_time") {
			// Skip the snapshot_time line
			continue
		}
		if strings.HasPrefix(line, "- job_id:") {
			suffix, _ := strings.CutPrefix(line, "- job_id:")
			job = strings.TrimSpace(suffix)
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 15 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		syscall := strings.TrimSuffix(fields[0], ":")

		numSamples, err := strconv.Atoi(strings.Trim(fields[3], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of samples: %v", err)
		}

		unit := strings.Trim(fields[5], ",")
		min, err := strconv.Atoi(strings.Trim(fields[7], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse min value: %v", err)
		}

		max, err := strconv.Atoi(strings.Trim(fields[9], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse max value: %v", err)
		}

		sum, err := strconv.Atoi(strings.Trim(fields[11], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum value: %v", err)
		}

		sumSquared, err := strconv.Atoi(strings.Trim(fields[13], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum squared value: %v", err)
		}

		result = append(result, JobStat{
			Job:        job,
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

type JobStatsCollector struct {
	jobStatsSamplesMetric *prometheus.Desc
	jobStatsSumMetric     *prometheus.Desc
	jobStatsSumsqMetric   *prometheus.Desc
	jobStatsFilePatterns  string
	level                 consts.Level
}

func NewJobStatsCollector(jobStatsSamplesMetric *prometheus.Desc, jobStatsSumMetric *prometheus.Desc, jobStatsSumsqMetric *prometheus.Desc, jobStatsFilePatterns string, level consts.Level) *JobStatsCollector {
	return &JobStatsCollector{
		jobStatsSamplesMetric: jobStatsSamplesMetric,
		jobStatsSumMetric:     jobStatsSumMetric,
		jobStatsSumsqMetric:   jobStatsSumsqMetric,
		jobStatsFilePatterns:  jobStatsFilePatterns,
		level:                 level,
	}
}

func (x *JobStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.jobStatsSamplesMetric
	ch <- x.jobStatsSumMetric
	ch <- x.jobStatsSumsqMetric
}

func (c *JobStatsCollector) CollectStatMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		value, err := os.ReadFile(filepath.Clean(path))
		if err != nil || value == nil {
			fmt.Printf("could not read stat file %s\n", path)
		}
		stats, err := ParseJobStat(string(value))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, stat := range stats {
			ch <- prometheus.MustNewConstMetric(c.jobStatsSamplesMetric, prometheus.GaugeValue, float64(stat.NumSamples), path, stat.Job, stat.Syscall)
			ch <- prometheus.MustNewConstMetric(c.jobStatsSumMetric, prometheus.GaugeValue, float64(stat.Sum), path, stat.Job, stat.Syscall, stat.Unit)
			ch <- prometheus.MustNewConstMetric(c.jobStatsSumsqMetric, prometheus.GaugeValue, float64(stat.SumSquared), path, stat.Job, stat.Syscall, stat.Unit)
		}
	}
}

// CollectBasicMetrics collects basic metrics
func (c *JobStatsCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Basic {
		c.CollectStatMetrics(ch, c.jobStatsFilePatterns)
	}
}

// CollectExtendedMetrics collects extended metrics
func (c *JobStatsCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	if c.level == consts.Extended {
		c.CollectStatMetrics(ch, c.jobStatsFilePatterns)
	}
}
