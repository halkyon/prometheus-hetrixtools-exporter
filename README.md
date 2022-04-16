# prometheus-hetrixtools-exporter

`prometheus-hetrixtools-exporter` exports [HetrixTools](https://hetrixtools.com) metrics to Prometheus.

Currently this exports two kinds of metrics.

An uptime status for each uptime monitor:

```
hetrixtools_uptime_monitor_status{id="12345",name="my check",port="123",status_text="Online",target="1.2.3.4"} 0
```

Response times from each location configured for each monitor:

```
hetrixtools_uptime_monitor_response_time{id="12345",location="Frankfurt",name="my check",port="123",target="1.2.3.4"} 99
hetrixtools_uptime_monitor_response_time{id="12345",location="New York",name="my check",port="123",target="1.2.3.4"} 18
hetrixtools_uptime_monitor_response_time{id="12345",location="San Francisco",name="my check",port="123",target="1.2.3.4"} 70
```

## Usage

`API_KEY` must be defined in the environment for this to work, which is your HetrixTools API key.

To run the exporter with Docker:

```
docker run -e API_KEY=mykey ghcr.io/halkyon/prometheus-hetrixtools-exporter
```

By default the exporter will listen on all interfaces on port `8080`. Use `-listen-address` to listen on something else.
