# prometheus-hetrixtools-exporter

`prometheus-hetrixtools-exporter` exports HetrixTools metrics to Prometheus.

Currently this exposes two kinds of metrics.

An uptime status for each check configured in HetrixTools:

```
hetrixtools_uptime_status{id="12345",name="my check",port="123",status_text="Online",target="1.2.3.4"} 0
```

Response times from each location configured for each check:

```
hetrixtools_uptime_monitor_response_time{id="12345",location="Frankfurt",name="my check",port="123",target="1.2.3.4"} 99
hetrixtools_uptime_monitor_response_time{id="12345",location="New York",name="my check",port="123",target="1.2.3.4"} 18
hetrixtools_uptime_monitor_response_time{id="12345",location="San Francisco",name="my check",port="123",target="1.2.3.4"} 70
```

## Usage

TODO
