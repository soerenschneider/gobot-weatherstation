# gobot-bme280
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-bme280)](https://goreportcard.com/report/github.com/soerenschneider/gobot-bme280)
![test-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/release.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/golangci-lint.yaml/badge.svg)

Detects and forwards temperature data using a [BME280 sensor](https://gobot.io/documentation/drivers/bme280/) and a Raspberry PI

## Features

🤖 Integrates with Home-Assistant<br/>
📊 Calculates statistics about read temperature data over time windows, accessible via MQTT and metrics<br/>
🔐 Allows connecting to secure MQTT brokers using TLS client certificates<br/>
🔭 Expose temperature data as metrics to enable alerting and Grafana dashboards<br/>

## Installation

### Binaries
Download a prebuilt binary from the [releases section](https://github.com/soerenschneider/gobot-bme280/releases) for your system.

### From Source
As a prerequisite, you need to have [Golang SDK](https://go.dev/dl/) installed. Then you can install gobot-bme280 from source by invoking:
```shell
$ go install github.com/soerenschneider/gobot-bme280@latest
```

## JSON Example Payload
```json
{"alt":99,"humidity":13,"pressure":13.37,"temp":22.25,"timestamp":1630563744}
```

## Configuration



## Metrics

This project exposes the following metrics in Open Metrics format using the `gobot_bme280` prefix.

| Metric Name                  | Description                                                       | Labels          |
|------------------------------|-------------------------------------------------------------------|-----------------|
| version                      | Version information of this robot                                 | version, commit |
| heartbeat_timestamp_seconds  | Heartbeat of this robot                                           | placement       |
| reading_errors_total         | Total amount of errors while reading from the sensor              | placement       |
| altitude_meters              | The measured altitude in meters                                   | placement       |
| humidity_percent             | The measured humidity in percent                                  | placement       |
| temperature_celsius          | The measured temperature in degrees celsius                       | placement       |
| pressure_pa                  | The measured pressure in pascal                                   | placement       |
| messages_published_total     | The amount of published MQTT messages                             | placement       |
| message_publish_errors_total | Total amount of errors while trying to publish messages over MQTT | placement       |
