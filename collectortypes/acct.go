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

type AcctCollector struct {
	AcctInodesMetric *prometheus.Desc
	AcctKbytesMetric *prometheus.Desc
	acctFilePatterns string
	acctFileRegex    regexp.Regexp
	BaseCollector
}

// SampleData represents the parsed information for each label
type AcctEntry struct {
	Id     string
	inodes int
	kbytes int
}

func NewAcctCollector(AcctInodesMetric *MetricInfo, AcctKbytesMetric *MetricInfo, acctFilePatterns string, acctFileRegex string, configKey string) *AcctCollector {
	acctFileRegexp := *regexp.MustCompile(acctFileRegex)
	return &AcctCollector{
		AcctInodesMetric: AcctInodesMetric.CreatePrometheusMetric([]string{"id"}, acctFileRegexp),
		AcctKbytesMetric: AcctKbytesMetric.CreatePrometheusMetric([]string{"id"}, acctFileRegexp),
		acctFilePatterns: acctFilePatterns,
		acctFileRegex:    acctFileRegexp,
		BaseCollector: BaseCollector{configKey: configKey},
	}
}

func (x *AcctCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.AcctInodesMetric
	ch <- x.AcctKbytesMetric
}

// CollectMetrics collects metrics
func (c *AcctCollector) CollectMetrics(ch chan<- prometheus.Metric) {
	c.CollectAcctMetrics(ch, c.acctFilePatterns)
}

func (c *AcctCollector) CollectAcctMetrics(ch chan<- prometheus.Metric, pattern string) {
	paths, _ := filepath.Glob(pattern)
	if paths == nil {
		return
	}
	for _, path := range paths {
		pathLabels := c.acctFileRegex.FindStringSubmatch(path)[1:]
		value, err := os.ReadFile(filepath.Clean(path))
		if err != nil || value == nil {
			fmt.Printf("could not read acct file %s\n", path)
		}
		accts, err := ParseAccts(string(value))
		if err != nil {
			fmt.Printf("got error while parsing line: %s\n", err)
		}
		for _, acct := range accts {
			ch <- prometheus.MustNewConstMetric(c.AcctInodesMetric, prometheus.GaugeValue, float64(acct.inodes), append([]string{acct.Id}, pathLabels...)...)
			ch <- prometheus.MustNewConstMetric(c.AcctKbytesMetric, prometheus.GaugeValue, float64(acct.kbytes), append([]string{acct.Id}, pathLabels...)...)
		}
	}
}

// ParseInput parses the input string and returns a slice of SampleData
func ParseAccts(input string) (map[string]AcctEntry, error) {
	var result map[string]AcctEntry = make(map[string]AcctEntry)

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
		if len(fields) < 7 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		inodes, err := strconv.Atoi(strings.Trim(fields[3], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse hard limit: %v", err)
		}

		kbytes, err := strconv.Atoi(strings.Trim(fields[5], ","))
		if err != nil {
			return nil, fmt.Errorf("failed to parse soft limit: %v", err)
		}

		result[id] = AcctEntry{
			Id:     id,
			inodes: inodes,
			kbytes: kbytes,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}
