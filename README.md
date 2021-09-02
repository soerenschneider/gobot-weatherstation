[![Go Report Card](https://goreportcard.com/badge/github.com/soerenschneider/gobot-weatherstation)](https://goreportcard.com/report/github.com/soerenschneider/gobot-weatherstation)

This project uses the [Gobot Framework](https://gobot.io/) in combination with a [BME280 sensor](https://gobot.io/documentation/drivers/bme280/) to work as a configurable weatherstation, being able to send JSON encoded payloads via MQTT and exposing machine readable metrics in the Open Metrics Format.

# JSON Example Payload
```json
{"alt":99,"humidity":13,"pressure":13.37,"temp":22.25,"timestamp":1630563744}
```

# Configuration
## Via Env Variables
| ENV                       | Default              | Description                                    |
|---------------------------|----------------------|------------------------------------------------|
| WEATHERBOT_LOCATION       | -                    | Location short name of this weatherstation bot |
| WEATHERBOT_I2C_BUS        | 1                    | I2C Bus to use                                 |
| WEATHERBOT_I2C_ADDRESS    | 0x76                 | I2C Address to use                             |
| WEATHERBOT_MQTT_HOST      | -                    | Host of the MQTT broker, can be omitted        |
| WEATHERBOT_MQTT_CLIENT_ID | weatherbot-$LOCATION | Client ID for the MQTT connection              |
| WEATHERBOT_MQTT_TOPIC     | weatherbot/$LOCATION | Topic to publish messages into                 |
| WEATHERBOT_METRICS_ADDR   | :9400                | Prometheus http handler listen address         |

## Via Config File

```json
{
  "location": "location",
  "read_interval": 60,
  "metrics_addr": ":1234",
  "i2c_bus": 15,
  "i2c_address": 16,
  "mqtt_host": "tcp://broker:1883",
  "mqtt_client_id": "client-id",
  "mqtt_topic": "mytopic/foo"
}
```

# Metrics

This project exposes the following metrics in Open Metrics format.

| Namespace  | Subsystem | Name                     | Type    | Labels   | Help                                                              |
|------------|-----------|--------------------------|---------|----------|-------------------------------------------------------------------|
| weatherbot | sensor    | reading_errors_total     | counter | location | Total amount of errors while reading from the sensor              |
| weatherbot | sensor    | altitude_meters          | gauge   | location | The measured altitude in meters                                   |
| weatherbot | sensor    | humidity_percent         | gauge   | location | The measured humidity in percent                                  |
| weatherbot | sensor    | temperature_celsius      | gauge   | location | The measured temperature in degrees celsius                       |
| weatherbot | sensor    | pressure_pa              | gauge   | location | The measured pressure in pascal                                   |
| weatherbot | mqtt      | messages_published_total | counter | location | The amount of published MQTT messages                             |
| weatherbot | mqtt      | message_publish_errors   | counter | location | Total amount of errors while trying to publish messages over MQTT |