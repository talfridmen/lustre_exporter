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
	// "github.com/talfridmen/lustre_exporter/consts"
)

type JobStatsCollector struct {
	jobStatsSamplesMetric *prometheus.Desc
	jobStatsSumMetric     *prometheus.Desc
	jobStatsFilePatterns  string
	jobStatsFileRegex     regexp.Regexp
	BaseCollector
}

// SampleData represents the parsed information for each label
type JobStat struct {
	Job                       string
	Syscall                   string
	NumSamples                int
	Unit                      string
	Min, Max, Sum, SumSquared int
}

func NewJobStatsCollector(jobStatsSamplesMetric *MetricInfo, jobStatsSumMetric *MetricInfo, jobStatsFilePatterns string, jobStatsFileRegex string, configName string) *JobStatsCollector {
	jobStatsFileRegexp := *regexp.MustCompile(jobStatsFileRegex)
	return &JobStatsCollector{
		jobStatsSamplesMetric: jobStatsSamplesMetric.CreatePrometheusMetric([]string{"job", "stat_type"}, jobStatsFileRegexp),
		jobStatsSumMetric:     jobStatsSumMetric.CreatePrometheusMetric([]string{"job", "stat_type", "units"}, jobStatsFileRegexp),
		jobStatsFilePatterns:  jobStatsFilePatterns,
		jobStatsFileRegex:     jobStatsFileRegexp,
		BaseCollector: BaseCollector{configKey: configName},
	}
}

func (x *JobStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.jobStatsSamplesMetric
	ch <- x.jobStatsSumMetric
}

// CollectMetrics collects metrics
func (c *JobStatsCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.jobStatsFilePatterns)
}

func (c *JobStatsCollector) CollectStatMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		pathLabels := c.jobStatsFileRegex.FindStringSubmatch(path)[1:]
		value, err := os.ReadFile(filepath.Clean(path))
		if err != nil || value == nil {
			fmt.Printf("could not read jobstat file %s\n", path)
		}
		stats, err := ParseJobStat(string(value))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, stat := range stats {
			ch <- prometheus.MustNewConstMetric(c.jobStatsSamplesMetric, prometheus.GaugeValue, float64(stat.NumSamples), append([]string{stat.Job, stat.Syscall}, pathLabels...)...)
			ch <- prometheus.MustNewConstMetric(c.jobStatsSumMetric, prometheus.GaugeValue, float64(stat.Sum), append([]string{stat.Job, stat.Syscall, stat.Unit}, pathLabels...)...)
		}
	}
}

type Key struct{ job, syscall string }

// ParseInput parses the input string and returns a slice of SampleData
func ParseJobStat(input string) (map[Key]JobStat, error) {
	var result map[Key]JobStat = map[Key]JobStat{}

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

		syscall := strings.TrimSuffix(fields[0], ":")

		numSamples, err := strconv.Atoi(strings.Trim(fields[3], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of samples: %v", err)
		}

		if numSamples == 0 {
			continue
		}

		unit := strings.Trim(fields[5], ",")

		if len(fields) < 12 {
			result[Key{job, syscall}] = JobStat{
				Job:        job,
				Syscall:    syscall,
				NumSamples: numSamples,
				Unit:       unit,
			}
			continue
		}

		sum, err := strconv.Atoi(strings.Trim(fields[11], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse sum value: %v", err)
		}

		result[Key{job, syscall}] = JobStat{
			Job:        job,
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
