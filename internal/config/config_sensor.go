package config

const (
	defaultGpioBus     = 1
	defaultGpioAddress = 0x76
)

func defaultSensorConfig() SensorConfig {
	return SensorConfig{
		GpioBus:     defaultGpioBus,
		GpioAddress: defaultGpioAddress,
	}
}

type SensorConfig struct {
	GpioBus     int `json:"gpio_bus,omitempty" env:"GPIO_BUS" validate:"gte=0"`
	GpioAddress int `json:"gpio_address,omitempty" env:"GPIO_ADDRESS" validate:"gte=1,lte=200"`
}
