package main

import (
	"flag"
	"fmt"
	"github.com/soerenschneider/gobot-bme280/internal"
	"github.com/soerenschneider/gobot-bme280/internal/config"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"os"
)

const (
	cliConfFile = "config"
	cliVersion  = "version"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, cliConfFile, "", "File to read configuration from")
	version := flag.Bool(cliVersion, false, "Print version and exit")

	flag.Parse()

	if *version {
		fmt.Printf("%s (revision %s)", internal.BuildVersion, internal.CommitHash)
		os.Exit(0)
	}

	log.Printf("Started %s, version %s, commit %s", config.BotName, internal.BuildVersion, internal.CommitHash)
	conf := loadConfig(configFile)
	conf.Print()
	log.Println("Validating config...")
	err := conf.Validate()
	if err != nil {
		log.Fatalf("Could not validate config: %v", err)
	}

	run(conf)
}

func run(conf config.Config) {
	if conf.MetricConfig != "" {
		go internal.StartMetricsServer(conf.MetricConfig)
	}

	log.Println("Building adaptors and drivers")
	raspberry := raspi.NewAdaptor()
	driver := i2c.NewBME280Driver(raspberry, i2c.WithBus(conf.GpioBus), i2c.WithAddress(conf.GpioAddress))

	var mqttAdaptor internal.WeatherBotMqttAdaptor
	if conf.Host != "" {
		log.Println("Building MQTT adaptor")

		clientId := fmt.Sprintf("%s_%s", config.BotName, conf.Placement)
		mq := mqtt.NewAdaptor(conf.MqttConfig.Host, clientId)
		mq.SetAutoReconnect(true)
		mq.SetQoS(1)

		if conf.MqttConfig.UsesSslCerts() {
			log.Println("Setting TLS client cert and key...")
			mq.SetClientCert(conf.MqttConfig.ClientCertFile)
			mq.SetClientKey(conf.MqttConfig.ClientKeyFile)

			if len(conf.MqttConfig.ServerCaFile) > 0 {
				log.Println("Setting server CA...")
				mq.SetServerCert(conf.MqttConfig.ServerCaFile)
			}
		}

		mqttAdaptor = mq
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
	err := bot.Start()
	if err != nil {
		log.Fatalf("Could not start bot: %v", err)
	}
}

func loadConfig(configFile string) config.Config {
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
