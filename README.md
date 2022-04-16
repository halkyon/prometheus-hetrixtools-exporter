# prometheus-hetrixtools-exporter

`prometheus-hetrixtools-exporter` exports [HetrixTools](https://hetrixtools.com) monitoring data to Prometheus.

Currently this exports two kinds of metrics.

An uptime status for each uptime monitor target:

```
hetrixtools_uptime_monitor_status{id="12345",name="my check",port="123",status_text="Online",target="1.2.3.4"} 0
```

TODO: This metric only accounts for "Online" status as being up, and anything else as down. It doesn't check for any
scheduled maintenance, for example.

Response times from each location monitoring a target:

```
hetrixtools_uptime_monitor_response_time_seconds{id="12345",location="Frankfurt",name="my check",port="123",target="1.2.3.4"} 0.099
hetrixtools_uptime_monitor_response_time_seconds{id="12345",location="New York",name="my check",port="123",target="1.2.3.4"} 0.018
hetrixtools_uptime_monitor_response_time_seconds{id="12345",location="San Francisco",name="my check",port="123",target="1.2.3.4"} 0.070
```

## Usage

`API_KEY` must be defined in the environment for this to work, which is your HetrixTools API key.

To run the exporter with Docker:

```
docker run -e API_KEY=mykey ghcr.io/halkyon/prometheus-hetrixtools-exporter
```

By default the exporter will listen on all interfaces on port `8080`. Use `-listen-address` to listen on something else.

Note that the HetrixTools free account is limited to 1000 API calls per month. If you intend to have Prometheus scrape
this data on a very frequent interval, you'll need to purchase a plan that includes more API calls.

## Building

[`goreleaser`](https://goreleaser.com) is used to automate the entire build and release process.

To run a development build, run `goreleaser --snapshot --rm-dist`.

To release a new version, create and push the tag in the format `v0.0.0`. Binaries and Docker images will be
built and released.
