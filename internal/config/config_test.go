package config

import (
	"reflect"
	"testing"
)

func Test_matchHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{
			name: "no tld",
			host: "tcp://hostname:1883",
			want: true,
		},
		{
			name: "tld",
			host: "tcp://hostname.my.tld:1883",
			want: true,
		},
		{
			name: "ip",
			host: "tcp://192.168.0.1:1883",
			want: true,
		},
		{
			name: "no protocol",
			host: "192.168.0.1:1883",
			want: false,
		},
		{
			name: "no port",
			host: "tcp://host",
			want: false,
		},
		{
			name: "only host",
			host: "host",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchHost(tt.host); got != tt.want {
				t.Errorf("matchHost() error = %v, got %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		placement    string
		MetricConfig string
		FirmAtaPort  string
		GpioBus      int
		GpioAddress  int
		IntervalSecs int
		LogValues    bool
		MqttConfig   MqttConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all okay",
			fields: fields{
				placement:    "placement",
				MetricConfig: "0.0.0.0:9100",
				FirmAtaPort:  "/dev/ttyUSB0",
				GpioBus:      1,
				GpioAddress:  75,
				IntervalSecs: 30,
				LogValues:    false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: false,
		},
		{
			name: "missing loc",
			fields: fields{
				MetricConfig: "0.0.0.0:9100",
				FirmAtaPort:  "/dev/ttyUSB0",
				GpioBus:      1,
				GpioAddress:  75,
				IntervalSecs: 30,
				LogValues:    false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid gpiobus",
			fields: fields{
				placement:    "loc",
				MetricConfig: "0.0.0.0:9100",
				GpioBus:      -5,
				GpioAddress:  75,
				IntervalSecs: 30,
				LogValues:    false,
				MqttConfig: MqttConfig{
					Host:  "tcp://host:80",
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
		{
			name: "missing host",
			fields: fields{
				placement:    "loc",
				MetricConfig: "0.0.0.0:9100",
				FirmAtaPort:  "/dev/ttyUSB0",
				GpioBus:      1,
				GpioAddress:  75,
				IntervalSecs: 30,
				LogValues:    false,
				MqttConfig: MqttConfig{
					Topic: "topic/bla",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Placement:    tt.fields.placement,
				MetricConfig: tt.fields.MetricConfig,
				SensorConfig: SensorConfig{
					GpioBus:     tt.fields.GpioBus,
					GpioAddress: tt.fields.GpioAddress,
				},
				IntervalSecs: tt.fields.IntervalSecs,
				LogSensor:    tt.fields.LogValues,
				MqttConfig:   tt.fields.MqttConfig,
			}
			if err := Validate(c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadJsonConfig(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     *Config
		wantErr  bool
	}{
		{
			name:     "non-existent-file",
			filePath: "ihopethispathdoesntexist/somefile.json",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "example-config",
			filePath: "../../contrib/example-config.json",
			want: &Config{
				Placement:    "location",
				MetricConfig: ":1234",
				SensorConfig: SensorConfig{
					GpioBus:     defaultGpioBus,
					GpioAddress: defaultGpioAddress,
				},
				IntervalSecs: defaultIntervalSeconds,
				LogSensor:    defaultLogSensor,
				MqttConfig: MqttConfig{
					Host:  "tcp://broker:1883",
					Topic: "mytopic/foo",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchTopic(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  bool
	}{
		{
			topic: "topicname",
			want:  true,
		},
		{
			topic: "more/complicated",
			want:  true,
		},
		{
			topic: "more/complicated/topic",
			want:  true,
		},
		{
			topic: "/leading",
			want:  false,
		},
		{
			topic: "trailing/",
			want:  false,
		},
		{
			topic: "replace/%s",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchTopic(tt.topic); got != tt.want {
				t.Errorf("matchTopic() got = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestConfig_TemplateTopic(t *testing.T) {
	type fields struct {
		placement  string
		MqttConfig MqttConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *Config
	}{
		{
			name: "template",
			fields: fields{
				placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/%s",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix/loc",
				},
			},
		},
		{
			name: "no templating",
			fields: fields{
				placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
			want: &Config{
				Placement: "loc",
				MqttConfig: MqttConfig{
					Topic: "prefix",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				Placement:  tt.fields.placement,
				MqttConfig: tt.fields.MqttConfig,
			}
			conf.FormatTopic()
			if !reflect.DeepEqual(conf, tt.want) {
				t.Fail()
			}
		})
	}
}
