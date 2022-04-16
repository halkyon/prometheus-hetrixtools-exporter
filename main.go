package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/halkyon/prometheus-hetrixtools-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const program = "hetrixtools_exporter"

func main() {
	displayVersion := flag.Bool("version", false, "Display version information")
	addr := flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	flag.Parse()

	versionCollector := version.NewCollector(program)

	if *displayVersion {
		println(version.Print(program))
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		log.Fatal("API_KEY is not defined in environment")
	}

	r := prometheus.NewRegistry()
	r.MustRegister(versionCollector)
	r.MustRegister(collector.New(apiKey))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	srv := http.Server{Addr: *addr, Handler: mux}
	log.Fatal(srv.ListenAndServe())
}
