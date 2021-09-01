package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const namespace = BotName

var (
	metricSensorErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "sensor_errors_total",
		Subsystem: "sensor",
		Help:      "The amount of errors while trying to read the sensor",
	}, []string{"location"})

	metricAltitude = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "altitude_meters",
		Subsystem: "sensor",
		Help:      "The measured altitude in meters",
	}, []string{"location"})

	metricHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "humidity_percent",
		Subsystem: "sensor",
		Help:      "The measured humidity in percent",
	}, []string{"location"})

	metricTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature_celsius",
		Subsystem: "sensor",
		Help:      "The measured temperature in degrees celsius",
	}, []string{"location"})

	metricPressure = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "pressure_pa",
		Subsystem: "sensor",
		Help:      "The measured pressure in pascal",
	}, []string{"location"})

	metricsMessagesPublished = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "messages_published_total",
		Subsystem: "mqtt",
		Help:      "The amount of published MQTT messages",
	}, []string{"location"})

	metricsMessagePublishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "message_publish_errors",
		Subsystem: "mqtt",
		Help:      "The amount of errors occurred while trying to publish messages over MQTT",
	}, []string{"location"})
)

func metricFromMeasurement(m Measurement, location string) {
	metricAltitude.WithLabelValues(location).Set(float64(m.Altitude))
	metricHumidity.WithLabelValues(location).Set(float64(m.Humidity))
	metricPressure.WithLabelValues(location).Set(float64(m.Pressure))
	metricTemperature.WithLabelValues(location).Set(float64(m.Temperature))
	if nil != m.Errors && len(m.Errors) > 0 {
		metricSensorErrors.WithLabelValues(location).Inc()
	}
}

func StartMetricsServer(listenAddr string) {
	log.Printf("Starting metrics listener at %s", listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("Could not start metrics listener: %v", err)
	}
}
