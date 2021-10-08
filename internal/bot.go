package internal

import (
	"gobot-bme280/internal/config"
	"gobot.io/x/gobot"
	"time"
)

type WeatherBotSensor interface {
	gobot.Driver
	Altitude() (alt float32, err error)
	Pressure() (press float32, err error)
	Temperature() (temp float32, err error)
	Humidity() (humidity float32, err error)
}

type WeatherBotMqttAdaptor interface {
	gobot.Connection
	Publish(topic string, msg []byte) bool
}

type WeatherBotAdaptors struct {
	Adaptor     gobot.Connection
	Driver      WeatherBotSensor
	MqttAdaptor WeatherBotMqttAdaptor
	Config      config.Config
}

func AssembleBot(bot *WeatherBotAdaptors) *gobot.Robot {
	versionInfo.WithLabelValues(BuildVersion, CommitHash).Set(1)
	work := func() {
		bot.readAndPublishMeasurement()
		gobot.Every(time.Duration(bot.Config.IntervalSecs)*time.Second, func() {
			bot.readAndPublishMeasurement()
			metricsHeartbeat.WithLabelValues(bot.Config.Location).SetToCurrentTime()
		})
	}

	adaptors := []gobot.Connection{bot.Adaptor}
	if bot.MqttAdaptor != nil {
		adaptors = append(adaptors, bot.MqttAdaptor)
	}
	robot := gobot.NewRobot(config.BotName,
		adaptors,
		[]gobot.Device{bot.Driver},
		work,
	)

	return robot
}

func (station *WeatherBotAdaptors) readAndPublishMeasurement() {
	measurement := station.readMeasurement()
	metricFromMeasurement(measurement, station.Config.Location)

	if station.MqttAdaptor != nil {
		msg, _ := measurement.AsJson()
		success := station.MqttAdaptor.Publish(station.Config.MqttConfig.Topic, msg)

		if success {
			metricsMessagesPublished.WithLabelValues(station.Config.Location).Inc()
		} else {
			metricsMessagePublishErrors.WithLabelValues(station.Config.Location).Inc()
		}
	}
}

func (station *WeatherBotAdaptors) readMeasurement() Measurement {
	measurement := NewMeasurement()
	measurement.AddAltitude(station.Driver.Altitude())
	measurement.AddHumidity(station.Driver.Humidity())
	measurement.AddPressure(station.Driver.Pressure())
	measurement.AddTemperature(station.Driver.Temperature())
	return measurement
}
