[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-weatherstation)](https://goreportcard.com/report/github.com/soerenschneider/gobot-weatherstation)

This project uses the [Gobot Framework](https://gobot.io/) in combination with a [BME280 sensor](https://gobot.io/documentation/drivers/bme280/) to work as a configurable weatherstation, being able to send JSON encoded payloads via MQTT and exposing machine readable metrics in the Open Metrics Format.

# JSON Example Payload
```json
{"alt":99,"humidity":13,"pressure":13.37,"temp":22.25,"timestamp":1630563744}
```

# Configuration
## Via Env Variables
| ENV                         | Default              | Description                                    |
|-----------------------------|----------------------|------------------------------------------------|
| GOBOT_BME280_LOCATION       | -                    | Location short name of this bot                |
| GOBOT_BME280_GPIO_BUS       | 1                    | GPIO Bus to use                                |
| GOBOT_BME280_GPIO_ADDRESS   | 0x76                 | GPIO Address to use                            |
| GOBOT_BME280_MQTT_HOST      | -                    | Host of the MQTT broker, can be omitted        |
| GOBOT_BME280_MQTT_TOPIC     |                      | Topic to publish messages into                 |
| GOBOT_BME280_LOG_SENSOR     | false                | Log sensor readings                            |
| GOBOT_BME280_METRICS_ADDR   | :9400                | Prometheus http handler listen address         |

## Via Config File

```json
{
  "location": "location",
  "interval_s": 60,
  "metrics_addr": ":1234",
  "gpio_bus": 15,
  "gpio_address": 16,
  "mqtt_host": "tcp://broker:1883",
  "mqtt_topic": "mytopic/%s",
  "log_sensor": true
}
```

# Metrics

This project exposes the following metrics in Open Metrics format.

| Namespace  | Subsystem | Name                     | Type    | Labels   | Help                                                              |
|------------|-----------|--------------------------|---------|----------|-------------------------------------------------------------------|
| gobot_bme280 | sensor    | reading_errors_total     | counter | location | Total amount of errors while reading from the sensor              |
| gobot_bme280 | sensor    | altitude_meters          | gauge   | location | The measured altitude in meters                                   |
| gobot_bme280 | sensor    | humidity_percent         | gauge   | location | The measured humidity in percent                                  |
| gobot_bme280 | sensor    | temperature_celsius      | gauge   | location | The measured temperature in degrees celsius                       |
| gobot_bme280 | sensor    | pressure_pa              | gauge   | location | The measured pressure in pascal                                   |
| gobot_bme280 | mqtt      | messages_published_total | counter | location | The amount of published MQTT messages                             |
| gobot_bme280 | mqtt      | message_publish_errors   | counter | location | Total amount of errors while trying to publish messages over MQTT |