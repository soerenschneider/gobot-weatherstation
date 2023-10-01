# gobot-bme280
[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-bme280)](https://goreportcard.com/report/github.com/soerenschneider/gobot-bme280)
![test-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/test.yaml/badge.svg)
![release-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/release.yaml/badge.svg)
![golangci-lint-workflow](https://github.com/soerenschneider/gobot-bme280/actions/workflows/golangci-lint.yaml/badge.svg)

Detects and forwards temperature data using a [BME280 sensor](https://gobot.io/documentation/drivers/bme280/) and a Raspberry PI

## Features

🤖 Integrates with Home-Assistant<br/>
📊 Reads sensor data, accessible via MQTT and metrics<br/>
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

gobot-bme280 can be fully configured using either environment variables or a config file. To supply a config file, the `-config` parameter is used.

### General Config Reference
| Struct Field      | Description                                  | Environment Variable              | Default Value   | Validation                               |
|-------------------|----------------------------------------------|-----------------------------------|-----------------|------------------------------------------|
| Placement         | Specifies the placement.                     | GOBOT_BME280_PLACEMENT            | N/A (required)  | required                                 |
| MetricConfig      | Metric server address.                       | GOBOT_BME280_METRICS_LISTEN_ADDR  | N/A (omitempty) | tcp_addr                                 |
| IntervalSecs      | Interval in seconds for sensor readings.     | GOBOT_BME280_INTERVAL_S           | 30              | min=30,max=300                           |
| StatIntervals     | Intervals for collecting statistics.         | GOBOT_BME280_STAT_INTERVALS       | N/A (dive)      | dive,min=10,max=3600                     |
| LogSensor         | Whether to log sensor readings.              | GOBOT_BME280_LOG_SENSOR_READINGS  | false           | N/A                                      |

### MQTT Config Reference
| Struct Field      | Description                               | Environment Variable                  | Default Value                                 | Validation                              |
|-------------------|-------------------------------------------|---------------------------------------|-----------------------------------------------|-----------------------------------------|
| Disabled          | Indicates if MQTT is disabled.            | GOBOT_BME280_MQTT_DISABLED            | false                                         | N/A                                     |
| Host              | MQTT broker host address.                 | GOBOT_BME280_MQTT_BROKER              | N/A (required_if=Disabled false, mqtt_broker) | required_if=Disabled false, mqtt_broker |
| Topic             | MQTT topic for sensor readings.           | GOBOT_BME280_MQTT_TOPIC               | N/A (required_if=Disabled false, mqtt_topic)  | required_if=Disabled false, mqtt_topic  |
| ClientKeyFile     | Client SSL key file for MQTT.             | GOBOT_BME280_MQTT_TLS_CLIENT_KEY_FILE | N/A (required_unless=ClientCertFile '', file) | required_unless=ClientCertFile '', file |
| ClientCertFile    | Client SSL certificate file for MQTT.     | GOBOT_BME280_MQTT_TLS_CLIENT_CRT_FILE | N/A (required_unless=ClientKeyFile '', file)  | required_unless=ClientKeyFile '', file  |
| ServerCaFile      | Server SSL CA certificate file for MQTT.  | GOBOT_BME280_MQTT_TLS_SERVER_CA_FILE  | N/A (omitempty, file)                         | required_unless=ClientKeyFile '', file  |

### Sensor Config Reference
| Struct Field      | Description               | Environment Variable          | Default Value | Validation      |
|-------------------|---------------------------|-------------------------------|---------------|-----------------|
| GpioBus           | GPIO bus for sensor.      | GOBOT_BME280_GPIO_BUS         | 1             | gte=0           |
| GpioAddress       | GPIO address for sensor.  | GOBOT_BME280_GPIO_ADDRESS     | 0x76          | gte=1,lte=200   |


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
