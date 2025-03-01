package collectortypes

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

type CollectorType interface {
	Describe(ch chan<- *prometheus.Desc)
	CollectMetrics(ch chan<- prometheus.Metric)
	GetConfigKey() string
}

type BaseCollector struct {
	configKey       string
}

type MetricInfo struct {
	name string
	help string
}

func NewMetricInfo(name string, help string) *MetricInfo {
	return &MetricInfo{
		name: name,
		help: help,
	}
}

func (m *MetricInfo) CreatePrometheusMetric(defaultLabels []string, pathReg regexp.Regexp) *prometheus.Desc {
	pathLabels := pathReg.SubexpNames()[1:]
	return prometheus.NewDesc(
		m.name,
		m.help,
		append(defaultLabels, pathLabels...),
		nil,
	)
}

func (c *BaseCollector) GetConfigKey() string {
	return c.configKey
}
