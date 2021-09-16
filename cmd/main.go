package main

import (
	"flag"
	"fmt"
	"gobot-bme280/internal"
	"gobot-bme280/internal/config"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

func main() {
	log.Printf("Started %s, version %s, commit %s, built at %s", config.BotName, internal.BuildVersion, internal.CommitHash, internal.BuildTime)
	conf := getConfig()
	log.Println("Validating config...")
	err := conf.Validate()
	conf.Print()
	if err != nil {
		log.Fatalf("Could not validate config: %v", err)
	}

	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	log.Println("Building adaptors and drivers")
	raspberry := raspi.NewAdaptor()
	driver := i2c.NewBME280Driver(raspberry, i2c.WithBus(conf.GpioBus), i2c.WithAddress(conf.GpioAddress))

	var mqttAdaptor internal.WeatherBotMqttAdaptor
	if conf.Host != "" {
		clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Location)
		mqttAdaptor = mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
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

func getConfig() config.Config {
	var configFile string
	flag.StringVar(&configFile, "config", "", "File to read configuration from")
	flag.Parse()
	if configFile == "" {
		log.Println("Building config from env vars")
		return config.ConfigFromEnv()
	}

	log.Printf("Reading config from file %s", configFile)
	conf, err := config.ReadJsonConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read config from %s: %v", configFile, err)
	}
	if nil == conf {
		log.Fatalf("Received empty config, should not happen")
	}
	return *conf
}
