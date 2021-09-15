package config

import (
	"fmt"
	"log"
)

const (
	defaultGpioBus              = 1
	defaultGpioAddress          = 0x76
)

func defaultSensorConfig() SensorConfig {
	return SensorConfig{
		GpioBus:     defaultGpioBus,
		GpioAddress: defaultGpioAddress,
	}
}

type SensorConfig struct {
	GpioBus     int 	`json:"gpio_bus,omitempty"`
	GpioAddress int    `json:"gpio_address,omitempty"`
}

func (conf *SensorConfig) Validate() error {
	if conf.GpioBus < 0 {
		return fmt.Errorf("invalid pin provided: %d", conf.GpioBus)
	}

	if conf.GpioAddress < 1 {
		return fmt.Errorf("polling interval must not be smaller than 5: %d", conf.GpioAddress)
	}

	if conf.GpioAddress > 200 {
		return fmt.Errorf("polling interval too high: %d", conf.GpioAddress)
	}

	return nil
}

func (conf *SensorConfig) ConfigFromEnv() {
	gpioBus, err := fromEnvInt("GPIO_BUS")
	if err == nil {
		conf.GpioBus = gpioBus
	}

	gpioAddress, err := fromEnvInt("GPIO_ADDRESS")
	if err == nil {
		conf.GpioAddress = gpioAddress
	}
}

func (conf *SensorConfig) Print() {
	log.Printf("GpioBus=%d", conf.GpioBus)
	log.Printf("GpioAddress=%d", conf.GpioAddress)
}
