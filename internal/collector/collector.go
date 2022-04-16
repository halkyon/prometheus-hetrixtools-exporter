package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	apiBaseURL     = "https://api.hetrixtools.com/v1"
	requestTimeout = 5 * time.Second
	namespace      = "hetrixtools"
)

// Collector is an implementation of Prometheus.Collector.
type Collector struct {
	apiKey              string
	client              http.Client
	totalScrapes        prometheus.Counter
	monitorUptimeStatus *prometheus.GaugeVec
	monitorResponseTime *prometheus.GaugeVec
	scrapeDurationTime  prometheus.Gauge
	errorDesc           *prometheus.Desc
}

type Monitor struct {
	ID            string            `json:"ID"`
	Name          string            `json:"Name"`
	Target        string            `json:"Target"`
	Port          string            `json:"Port"`
	UptimeStatus  string            `json:"Uptime_Status"`
	ResponseTimes map[string]string `json:"Response_Time"`
}

// New returns a new Collector.
func New(apiKey string) *Collector {
	totalScrapes := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "scrapes_total",
		Help:      "Total scrapes",
	})

	monitorUptimeStatus := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "uptime_monitor_status",
		Help:      "Uptime status of recent monitor check (0: up, 1: down)",
	}, []string{"id", "name", "target", "port", "status_text"})

	monitorResponseTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "uptime_monitor_response_time",
		Help:      "Response time of recent monitor check in milliseconds",
	}, []string{"id", "name", "location", "target", "port"})

	scrapeDurationTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "scrape_duration_seconds",
		Help:      "Time taken to scrape metrics",
	})

	errorDesc := prometheus.NewDesc(fmt.Sprintf("%s_scrape_error", namespace), "Error scraping target", nil, nil)

	return &Collector{
		apiKey: apiKey,
		client: http.Client{
			Timeout: requestTimeout,
		},
		totalScrapes:        totalScrapes,
		monitorUptimeStatus: monitorUptimeStatus,
		monitorResponseTime: monitorResponseTime,
		scrapeDurationTime:  scrapeDurationTime,
		errorDesc:           errorDesc,
	}
}

// Describe implements Prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.totalScrapes.Describe(ch)
	c.monitorUptimeStatus.Describe(ch)
	c.monitorResponseTime.Describe(ch)
	c.scrapeDurationTime.Describe(ch)
}

// Collect implements Prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	c.totalScrapes.Inc()
	c.totalScrapes.Collect(ch)

	c.collectUptimeMonitors(ch)

	c.scrapeDurationTime.Set(time.Since(start).Seconds())
	c.scrapeDurationTime.Collect(ch)
}

func (c *Collector) collectUptimeMonitors(ch chan<- prometheus.Metric) {
	monitors, err := fetchMonitors(c.client, c.apiKey, "uptime/monitors/0/5000")
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.errorDesc, err)
		return
	}

	for _, mon := range monitors {
		var status float64
		// todo: what other values does UptimeStatus contain?
		if mon.UptimeStatus != "Online" {
			status = 1
		}

		c.monitorUptimeStatus.WithLabelValues(mon.ID, mon.Name, mon.Target, mon.Port, mon.UptimeStatus).Set(status)

		for loc, time := range mon.ResponseTimes {
			loc = strings.ReplaceAll(loc, "_", " ")

			timeFloat, err := strconv.ParseFloat(time, 64)
			if err != nil {
				ch <- prometheus.NewInvalidMetric(c.errorDesc, err)
			}

			c.monitorResponseTime.WithLabelValues(mon.ID, mon.Name, loc, mon.Target, mon.Port).Set(timeFloat)
		}
	}

	c.monitorUptimeStatus.Collect(ch)
	c.monitorResponseTime.Collect(ch)
}

func fetchMonitors(client http.Client, apiKey, endpoint string) (monitors []Monitor, err error) {
	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return monitors, fmt.Errorf("invalid url: %s", apiBaseURL)
	}
	u.Path = path.Join(u.Path, apiKey, endpoint)

	req, err := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return monitors, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		return monitors, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return monitors, fmt.Errorf("bad response: %d", rsp.StatusCode)
	}

	var tmp []json.RawMessage
	if err := json.NewDecoder(rsp.Body).Decode(&tmp); err != nil {
		return monitors, fmt.Errorf("json decoder: %w", err)
	}

	if err := json.Unmarshal(tmp[0], &monitors); err != nil {
		return monitors, fmt.Errorf("json unmarshal checks: %w", err)
	}

	return monitors, nil
}
