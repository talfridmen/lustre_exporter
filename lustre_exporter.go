package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"os"

	"gopkg.in/ini.v1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/talfridmen/lustre_exporter/collectors"
)

func main() {
	port := flag.Int("port", 9090, "Port to expose metrics on (defaults to 9090)")
	config := flag.String("config", "/etc/lustre_exporter.ini", "configuration file path")
	flag.Parse()

	inidata, err := ini.Load(*config)
	if err != nil {
		fmt.Printf("Fail to read ini file: %v", err)
		os.Exit(1)
	}

	// Create a new exporter
	exporter := NewExporter()

	// Register collectors with user-specified levels
	exporter.RegisterCollector(collectors.NewMDTCollector("mdt", inidata.Section("mdt")))
	exporter.RegisterCollector(collectors.NewOBDFilterCollector("obdfilter", inidata.Section("obdfilter")))
	exporter.RegisterCollector(collectors.NewOsdCollector("osd", inidata.Section("osd")))
	exporter.RegisterCollector(collectors.NewLdlmCollector("ldlm", inidata.Section("ldlm")))
	exporter.RegisterCollector(collectors.NewLliteCollector("llite", inidata.Section("client")))

	// Start the exporter
	exporter.Start(fmt.Sprintf(`:%d`, *port))
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

	fmt.Printf("Registering collector %s... ", c.GetName())
	e.collectors = append(e.collectors, c)
	fmt.Printf("Success!\n")
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(e.collectors))
	for _, c := range e.collectors {
		go func(collector collectors.Collector) {
			defer wg.Done()
			collector.CollectMetrics(ch)
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
