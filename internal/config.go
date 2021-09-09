package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const BotName = "weatherbot"

// This regex is not a very strict check, we don't validate hostname or ip (v4, v6) addresses...
var mqttHostRegex = regexp.MustCompile(`\w{3,}://.{3,}:\d{2,4}`)

type Config struct {
	Location        string `json:"location"`
	IntervalSeconds int    `json:"read_interval"`
	MetricConfig    string `json:"metrics_addr"`
	MqttConfig
	I2cConfig
}

type I2cConfig struct {
	Bus     int `json:"i2c_bus"`
	Address int `json:"i2c_address"`
}

type MqttConfig struct {
	Host     string `json:"mqtt_host"`
	ClientId string `json:"mqtt_client_id"`
	Topic    string `json:"mqtt_topic"`
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

func ReadJsonConfig(filePath string) (*Config, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config from file: %v", err)
	}

	ret := &Config{}
	err = json.Unmarshal(fileContent, ret)
	return ret, err
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

	return matchHost(c.MqttConfig.Host)
}

func matchHost(host string) error {
	if !mqttHostRegex.Match([]byte(host)) {
		return fmt.Errorf("invalid host format used")
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

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}
