package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
)

const (
	BotName                = "gobot_bme280"
	defaultLogSensor       = false
	defaultIntervalSeconds = 30
	defaultMetricConfig    = ":9192"
)

var (
	once     sync.Once
	validate *validator.Validate
)

type Config struct {
	Placement    string `json:"placement,omitempty"`
	MetricConfig string `json:"metrics_addr,omitempty"`
	IntervalSecs int    `json:"interval_s,omitempty" validate:"gte=30,lte=300"`
	LogSensor    bool   `json:"log_sensor,omitempty"`
	DisableMqtt  bool   `json:"disable_mqtt"`
	MqttConfig
	SensorConfig
}

func DefaultConfig() Config {
	return Config{
		LogSensor:    defaultLogSensor,
		IntervalSecs: defaultIntervalSeconds,
		MetricConfig: defaultMetricConfig,
		SensorConfig: defaultSensorConfig(),
	}
}

func Read(filePath string) (*Config, error) {
	ret := DefaultConfig()

	if len(filePath) > 0 {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not read config from file: %v", err)
		}

		err = json.Unmarshal(fileContent, &ret)
		if err != nil {
			return nil, err
		}
	}

	opts := env.Options{
		Prefix: fmt.Sprintf("%s_", strings.ToUpper(BotName)),
	}
	err := env.ParseWithOptions(&ret, opts)
	return &ret, err
}

func Validate(s interface{}) error {
	once.Do(func() {
		validate = validator.New()
		if err := validate.RegisterValidation("mqtt_topic", validateTopic); err != nil {
			log.Fatal("could not build custom validation 'mqtt_topic'")
		}
		if err := validate.RegisterValidation("mqtt_broker", validateBroker); err != nil {
			log.Fatal("could not build custom validation 'validateBroker'")
		}
	})
	return validate.Struct(s)
}

func validateTopic(fl validator.FieldLevel) bool {
	// Get the field value and check if it's a slice
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	topic, ok := field.Interface().(string)
	if !ok || !matchTopic(topic) {
		return false
	}

	return true
}

func validateBroker(fl validator.FieldLevel) bool {
	// Get the field value and check if it's a slice
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	// Convert to string and check its value
	broker, ok := field.Interface().(string)
	if !ok || !matchHost(broker) {
		return false
	}

	return true
}
