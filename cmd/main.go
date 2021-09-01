package main

import (
	"gobot-weatherstation/internal"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

func main() {
	conf := internal.DefaultConfig()
	err := conf.Validate()
	if err != nil {
		log.Fatalf("Could not validate config: %v", err)
	}

	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	raspberry := raspi.NewAdaptor()
	driver := i2c.NewBME280Driver(raspberry, i2c.WithBus(conf.I2cConfig.Bus), i2c.WithAddress(conf.I2cConfig.Address))

	var mqttAdaptor internal.WeatherBotMqttAdaptor
	if conf.Host != "" {
		mqttAdaptor = mqtt.NewAdaptor(conf.MqttConfig.Host, conf.MqttConfig.ClientId)
	} else {
		log.Println("No MQTT host defined, not connecting to MQTT broker")
	}

	adaptors := &internal.WeatherBotAdaptors{
		Driver:      driver,
		Adaptor:     raspberry,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	bot := internal.AssembleBot(adaptors)
	err = bot.Start()
	if err != nil {
		log.Fatalf("Could not start bot: %v", err)
	}
}
