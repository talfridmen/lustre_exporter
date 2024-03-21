package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/talfridmen/lustre_exporter/collectors"
	"github.com/talfridmen/lustre_exporter/consts"
)

func main() {
	// Define command-line flags
	mdtLevel := flag.String("mdt", "disabled", "Enable mdt collection (disabled,basic,extended)")
	obdfilterLevel := flag.String("obdfilter", "disabled", "Enable obdfilter collection (disabled,basic,extended)")
	osdLevel := flag.String("osd", "disabled", "Enable osd collection (disabled,basic,extended)")
	ldlmLevel := flag.String("ldlm", "disabled", "Enable ldlm collection (disabled,basic,extended)")
	clientLevel := flag.String("client", "disabled", "Enable client collection (disabled,basic,extended)")
	// Parse command-line flags
	flag.Parse()

	// Create a new exporter
	exporter := NewExporter()

	// Register collectors with user-specified levels
	exporter.RegisterCollector(collectors.NewMDTCollector("mdt", *mdtLevel))
	exporter.RegisterCollector(collectors.NewOBDFilterCollector("obdfilter", *obdfilterLevel))
	exporter.RegisterCollector(collectors.NewOsdCollector("osd", *osdLevel))
	exporter.RegisterCollector(collectors.NewLdlmCollector("ldlm", *ldlmLevel))
	exporter.RegisterCollector(collectors.NewLliteCollector("client", *clientLevel))

	// Start the exporter
	exporter.Start(":9090")
}

// Exporter represents the Prometheus exporter
type Exporter struct {
	mu         sync.Mutex
	collectors []collectors.Collector
}

// NewExporter creates a new exporter instance
func NewExporter() *Exporter {
	return &Exporter{}
}

// RegisterCollector registers a collector to the exporter
func (e *Exporter) RegisterCollector(c collectors.Collector) {
	e.mu.Lock()
	defer e.mu.Unlock()

	fmt.Printf("Registering collector %s level %d...", c.GetName(), c.GetLevel())
	e.collectors = append(e.collectors, c)
	fmt.Printf("Success!\n")
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(e.collectors))
	for _, c := range e.collectors {
		go func(collector collectors.Collector) {
			defer wg.Done()
			switch collector.GetLevel() {
			case consts.Disabled:
			case consts.Basic:
				collector.CollectBasicMetrics(ch)
			case consts.Extended:
				collector.CollectBasicMetrics(ch)
				collector.CollectExtendedMetrics(ch)
			}
		}(c)
	}
	wg.Wait()
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range e.collectors {
		c.Describe(ch)
	}
}

// Start starts the exporter on the given address
func (e *Exporter) Start(address string) {
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(address, nil))
}
