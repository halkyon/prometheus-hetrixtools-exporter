package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/halkyon/prometheus-hetrixtools-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promversion "github.com/prometheus/common/version"
)

var (
	version = "dev"
	commit  = ""
	date    = ""

	displayVersion = flag.Bool("version", false, "Display version information")
	addr           = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	versionGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": version,
		},
	})
)

func main() {
	flag.Parse()

	if *displayVersion {
		printVersion()
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		log.Fatal("API_KEY is not defined in environment")
	}

	prometheus.MustRegister(promversion.NewCollector("hetrixtools_exporter"))

	r := prometheus.NewRegistry()
	r.MustRegister(collector.UptimeStatus)
	r.MustRegister(collector.ResponseTime)
	r.MustRegister(versionGauge)

	c := collector.New(apiKey)
	r.MustRegister(c)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	srv := http.Server{Addr: *addr, Handler: mux}
	log.Fatal(srv.ListenAndServe())
}

func printVersion() {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}

	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}

	fmt.Println(result)
}
