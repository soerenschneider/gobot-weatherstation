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
	GpioBus     int `json:"gpio_bus,omitempty" validate:"gte=0"`
	GpioAddress int `json:"gpio_address,omitempty" validate:"gte=1,lte=200"`
}
