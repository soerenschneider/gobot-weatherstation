package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soerenschneider/gobot-bme280/internal/config"
	"log"
	"net/http"
)

const namespace = config.BotName

var (
	versionInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "version",
		Help:      "Version information of this robot",
	}, []string{"version", "commit"})

	metricsHeartbeat = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "heartbeat_timestamp_seconds",
		Help:      "Heartbeat of this robot",
	}, []string{"location"})

	metricSensorErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "reading_errors_total",
		Subsystem: "sensor",
		Help:      "Total amount of errors while reading from the sensor",
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
		Name:      "message_publish_errors_total",
		Subsystem: "mqtt",
		Help:      "Total amount of errors while trying to publish messages over MQTT",
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
