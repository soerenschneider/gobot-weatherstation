package internal

import (
	"reflect"
	"testing"
)

func TestReadJsonConfig(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name:    "non-existing conf",
			args:    args{filePath: "yada-dada-doo.json"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid config",
			args: args{filePath: "../contrib/test.json"},
			want: &Config{
				Location:        "location",
				IntervalSeconds: 60,
				MetricConfig:    ":1234",
				MqttConfig: MqttConfig{
					Host:     "tcp://broker:1883",
					ClientId: "client-id",
					Topic:    "mytopic/foo",
				},
				I2cConfig: I2cConfig{
					Bus:     15,
					Address: 16,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadJsonConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJsonConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJsonConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{
			name: "valid host",
			host: "tcp://myhost:1883",
			wantErr: false,
		},
		{
			name: "valid ip",
			host: "tcp://192.168.1.1:1883",
			wantErr: false,
		},
		{
			name: "missing protocol",
			host: "192.168.1.1:1883",
			wantErr: true,
		},
		{
			name: "missing port",
			host: "tcp://192.168.1.1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchHost(tt.host); (err != nil) != tt.wantErr {
				t.Errorf("matchHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}