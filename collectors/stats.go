package collectors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO: remove after use
// var (
// 	basic_stats_file_patterns    = [...]string{"mdt/*/md_stats", "obdfilter/*/stats"}
// 	extended_stats_file_patterns = [...]string{"mdt/*/exports/*/stats", "obdfilter/*/exports/*/stats", "ldlm.namespaces.filter-*.pool.stats"}
// )

type statsCollector struct {
	stats_samples_metric         *prometheus.Desc
	stats_sum_metric             *prometheus.Desc
	stats_sumsq_metric           *prometheus.Desc
	basic_stats_file_patterns    []string
	extended_stats_file_patterns []string
}

func (x *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.stats_samples_metric
	ch <- x.stats_sum_metric
	ch <- x.stats_sumsq_metric
}

func (c *statsCollector) CollectStatMetrics(ch chan<- prometheus.Metric, patterns []string) {
	for _, pattern := range patterns {
		fullPath := filepath.Join(procfs_base_dir, pattern)
		paths, _ := filepath.Glob(fullPath)
		if paths == nil {
			continue
		}
		for _, path := range paths {
			patternPath := strings.Replace(path, procfs_base_dir, "", 1)
			value, err := os.ReadFile(filepath.Clean(path))
			if err != nil || value == nil {
				fmt.Printf("could not read stat file %s\n", path)
			}
			stats, err := ParseStats(string(value))
			if err != nil {
				fmt.Printf("got error while parsing line: %s\n", err)
			}
			for _, stat := range stats {
				ch <- prometheus.MustNewConstMetric(c.stats_samples_metric, prometheus.GaugeValue, float64(stat.NumSamples), patternPath, stat.Syscall)
				ch <- prometheus.MustNewConstMetric(c.stats_sum_metric, prometheus.GaugeValue, float64(stat.Sum), patternPath, stat.Syscall, stat.Unit)
				ch <- prometheus.MustNewConstMetric(c.stats_sumsq_metric, prometheus.GaugeValue, float64(stat.SumSquared), patternPath, stat.Syscall, stat.Unit)
			}
		}
	}
}

// CollectBasicMetrics collects basic metrics
func (c *statsCollector) CollectBasicMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.basic_stats_file_patterns[:])
}

// CollectExtendedMetrics collects extended metrics
func (c *statsCollector) CollectExtendedMetrics(ch chan<- prometheus.Metric) {
	c.CollectStatMetrics(ch, c.extended_stats_file_patterns[:])
}
