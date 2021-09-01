package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const BotName = "weatherbot"

type Config struct {
	Location        string
	IntervalSeconds int
	MetricConfig    string
	MqttConfig
	I2cConfig
}

type I2cConfig struct {
	Bus     int
	Address int
}

type MqttConfig struct {
	Host     string
	ClientId string
	Topic    string
}

func DefaultConfig() Config {
	location := fromEnv(fmt.Sprintf("%s_LOCATION", strings.ToUpper(BotName)), "")
	return Config{
		Location:        location,
		IntervalSeconds: 60,
		I2cConfig: I2cConfig{
			Bus:     fromEnvInt(fmt.Sprintf("%s_I2C_BUS", strings.ToUpper(BotName)), 1),
			Address: fromEnvInt(fmt.Sprintf("%s_I2C_ADDRESS", strings.ToUpper(BotName)), 0x76),
		},
		MqttConfig: MqttConfig{
			Host:     fromEnv(fmt.Sprintf("%s_MQTT_HOST", strings.ToUpper(BotName)), ""),
			ClientId: fromEnv(fmt.Sprintf("%s_MQTT_CLIENT_ID", strings.ToUpper(BotName)), fmt.Sprintf("%s-%s", BotName, location)),
			Topic:    fromEnv(fmt.Sprintf("%s_MQTT_TOPIC", strings.ToUpper(BotName)), fmt.Sprintf("%s/%s", BotName, location)),
		},
		MetricConfig: fromEnv(fmt.Sprintf("%s_METRICS_ADDR", strings.ToUpper(BotName)), ":9400"),
	}
}

func (c *Config) Validate() error {
	if len(c.Location) == 0 {
		return fmt.Errorf("empty location provided")
	}

	if c.I2cConfig.Address < 0 {
		return fmt.Errorf("invalid i2c address provided")
	}

	if c.I2cConfig.Bus < 0 {
		return fmt.Errorf("invalid i2c bus provided")
	}

	// TODO: improve check
	if strings.Index(c.MqttConfig.Topic, " ") != -1 {
		return fmt.Errorf("invalid mqtt topic provided")
	}

	return nil
}

func fromEnv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}
	return val
}

func fromEnvInt(name string, def int) int {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	parsed, err := strconv.Atoi(name)
	if err != nil {
		return def
	}
	return parsed
}
