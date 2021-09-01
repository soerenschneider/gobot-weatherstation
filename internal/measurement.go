package internal

import (
	"encoding/json"
	"log"
	"time"
)

type Measurement struct {
	Altitude    float32  `json:"alt"`
	Humidity    float32  `json:"humidity"`
	Pressure    float32  `json:"pressure"`
	Temperature float32  `json:"temp"`
	Timestamp   int64    `json:"timestamp"`
	Errors      []string `json:"errors,omitempty"`
}

func NewMeasurement() *Measurement {
	return &Measurement{
		Altitude:    -1,
		Humidity:    -1,
		Pressure:    -1,
		Temperature: -1,
		Timestamp:   time.Now().Unix(),
		Errors:      nil,
	}
}

func (m Measurement) AsJson() ([]byte, error) {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Println("Error marshalling json")
	}
	return msg, err
}

func (m *Measurement) AddAltitude(alt float32, err error) {
	if err != nil {
		if m.Errors == nil {
			m.Errors = make([]string, 4)
		}
		m.Errors = append(m.Errors, err.Error())
		log.Printf("Error occurred while reading altitude from sensor: %v", err)
	} else {
		m.Altitude = alt
	}
}

func (m *Measurement) AddHumidity(hum float32, err error) {
	if err != nil {
		if m.Errors == nil {
			m.Errors = make([]string, 4)
		}
		m.Errors = append(m.Errors, err.Error())
		log.Printf("Error while reading humidity from sensor: %v", err)
	} else {
		m.Humidity = hum
	}
}

func (m *Measurement) AddPressure(pressure float32, err error) {
	if err != nil {
		if m.Errors == nil {
			m.Errors = make([]string, 4)
		}
		m.Errors = append(m.Errors, err.Error())
		log.Printf("Error while reading pressure from sensor: %v", err)
	} else {
		m.Pressure = pressure
	}
}

func (m *Measurement) AddTemperature(temp float32, err error) {
	if err != nil {
		if m.Errors == nil {
			m.Errors = make([]string, 4)
		}
		m.Errors = append(m.Errors, err.Error())
		log.Printf("Error while reading temperature from sensor: %v", err)
	} else {
		m.Temperature = temp
	}
}
