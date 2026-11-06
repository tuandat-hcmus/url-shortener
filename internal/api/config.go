package api

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark-service"`
	InstanceID  string `envconfig:"INSTANCE_ID"`
	Port        int    `envconfig:"PORT" default:"8080"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
