package internal

import (
	"encoding/json"
	"github.com/soerenschneider/gobot-bme280/internal/config"
	"gobot.io/x/gobot"
	"log"
	"testing"
	"time"
)

const (
	MeasureDefaultsPressure    = 13.37
	MeasureDefaultsHumidity    = 13.0
	MeasureDefaultsTemperature = 22.25
	MeasureDefaultsAltitude    = 99.0
)

func TestAssembleBot(t *testing.T) {
	conf := config.DefaultConfig()
	mqttAdaptor := &FakeMqttAdapter{}
	fakeAdaptor := &FakeMqttAdapter{} // this adaptor isn't really being used by our fake adaptors
	station := &WeatherBotAdaptors{
		Driver:      &FakeBme280{Conn: fakeAdaptor},
		Adaptor:     fakeAdaptor,
		MqttAdaptor: mqttAdaptor,
		Config:      conf,
	}

	// check preconditions
	if mqttAdaptor.Topic != "" {
		t.Errorf("Failed precondition, topic isn't empty: %s", mqttAdaptor.Topic)
	}
	if len(mqttAdaptor.Msg) != 0 {
		t.Errorf("Failed precondition, msg isn't empty: %s", mqttAdaptor.Msg)
	}

	bot := AssembleBot(station)
	go func() {
		_ = bot.Start()
	}()

	// TODO: Come on, man, fix this
	time.Sleep(5 * time.Second)
	if err := bot.Stop(); err != nil {
		t.Fatal(err)
	}

	m := &Measurement{}
	if err := json.Unmarshal(mqttAdaptor.Msg, m); err != nil {
		t.Fatal(err)
	}

	if m.Pressure != MeasureDefaultsPressure {
		t.Errorf("Expected %f, got %f", MeasureDefaultsPressure, m.Pressure)
	}

	if m.Temperature != MeasureDefaultsTemperature {
		t.Errorf("Expected %f, got %f", MeasureDefaultsTemperature, m.Temperature)
	}

	if m.Humidity != MeasureDefaultsHumidity {
		t.Errorf("Expected %f, got %f", MeasureDefaultsHumidity, m.Humidity)
	}

	if m.Altitude != MeasureDefaultsAltitude {
		t.Errorf("Expected %f, got %f", MeasureDefaultsAltitude, m.Altitude)
	}
}

type FakeMqttAdapter struct {
	Msg   []byte
	Topic string
}

func (m *FakeMqttAdapter) Name() string {
	return "FakeMqttAdapter"
}
func (m *FakeMqttAdapter) SetName(n string) {

}
func (m *FakeMqttAdapter) Connect() error {
	return nil
}
func (m *FakeMqttAdapter) Finalize() error {
	return nil
}

func (m *FakeMqttAdapter) Publish(topic string, msg []byte) bool {
	m.Topic = topic
	m.Msg = msg
	log.Printf("%s -> %v", topic, string(msg))
	return true
}

// ---------------------

type FakeBme280 struct {
	Conn gobot.Connection
}

func (driver *FakeBme280) Name() string {
	return "FakeBme280"
}
func (driver *FakeBme280) SetName(s string) {}

func (driver *FakeBme280) Start() error {
	return nil
}
func (driver *FakeBme280) Halt() error {
	return nil
}
func (driver *FakeBme280) Connection() gobot.Connection {
	return driver.Conn
}

func (driver *FakeBme280) Altitude() (alt float32, err error) {
	return MeasureDefaultsAltitude, nil
}

func (driver *FakeBme280) Pressure() (press float32, err error) {
	return MeasureDefaultsPressure, nil
}

func (driver *FakeBme280) Temperature() (temp float32, err error) {
	return MeasureDefaultsTemperature, nil
}

func (driver *FakeBme280) Humidity() (humidity float32, err error) {
	return MeasureDefaultsHumidity, nil
}
