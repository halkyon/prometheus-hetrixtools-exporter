package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	apiBaseURL     = "https://api.hetrixtools.com/v1"
	requestTimeout = 5 * time.Second
)

var (
	UptimeStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hetrixtools_uptime_status",
		Help: "Uptime status of recent check (0: up, 1: down)",
	}, []string{"id", "name", "target", "port", "status_text"})

	ResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hetrixtools_uptime_monitor_response_time",
		Help: "Response time of recent check in milliseconds",
	}, []string{"id", "name", "location", "target", "port"})
)

// Collector is an implementation of Prometheus.Collector.
type Collector struct {
	apiKey string
	client http.Client
}

type checks []struct {
	ID            string            `json:"ID"`
	Name          string            `json:"Name"`
	Target        string            `json:"Target"`
	Port          string            `json:"Port"`
	UptimeStatus  string            `json:"Uptime_Status"`
	ResponseTimes map[string]string `json:"Response_Time"`
}

// New returns a new Collector.
func New(apiKey string) *Collector {
	return &Collector{
		apiKey: apiKey,
		client: http.Client{
			Timeout: requestTimeout,
		},
	}
}

// Describe implements Prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect implements Prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	var checks checks
	if err := fetchChecks(c.client, c.apiKey, "uptime/monitors/0/30", &checks); err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("hetrixtools_scrape_error", "Error scraping target", nil, nil), err)
		return
	}

	for _, check := range checks {
		var status float64
		// todo: what other values does UptimeStatus contain?
		if check.UptimeStatus != "Online" {
			status = 1
		}

		UptimeStatus.WithLabelValues(
			check.ID,
			check.Name,
			check.Target,
			check.Port,
			check.UptimeStatus,
		).Set(status)

		for location, time := range check.ResponseTimes {
			location = strings.ReplaceAll(location, "_", " ")

			timeFloat, err := strconv.ParseFloat(time, 64)
			if err != nil {
				ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("hetrixtools_scrape_error", "Error scraping target", nil, nil), err)
			}

			ResponseTime.WithLabelValues(
				check.ID,
				check.Name,
				location,
				check.Target,
				check.Port,
			).Set(timeFloat)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("hetrixtools_scrape_duration_seconds", "Total HetrixTools scrape time", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())
}

func fetchChecks(client http.Client, apiKey, endpoint string, c *checks) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", apiBaseURL, apiKey, endpoint), http.NoBody)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	var tmp []json.RawMessage
	if err := json.NewDecoder(rsp.Body).Decode(&tmp); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	if err := json.Unmarshal(tmp[0], c); err != nil {
		return fmt.Errorf("json unmarshal checks: %w", err)
	}

	return nil
}
