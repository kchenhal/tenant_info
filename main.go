package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	// promethus
	"github.com/kchenhal/client_golang/prometheus"
	"github.com/kchenhal/client_golang/prometheus/promhttp"
)

var (
	addr     = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	logLevel = flag.Int64("log-level", 4, "value from 1 to 5, the higher number, less output")
)

var (
	tenantInterval = time.Duration(60) * time.Minute
	r1             = rand.New(rand.NewSource(time.Now().UnixNano()))

	tenantUsage = prometheus.NewGaugeVecReset(
		prometheus.GaugeOpts{
			Name: "tenant_usage",
			Help: "Tenant last use data",
		},
		[]string{"dc", "name", "query_date"},
	)
)

func init() {
	prometheus.MustRegister(tenantUsage)
}

func main() {
	flag.Parse()

	go generateTenantInfo()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: false,
		},
	))
	http.HandleFunc("/logLevel", handleLogLevel)
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/tenantInterval", handleUpdateTenantInterval)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
