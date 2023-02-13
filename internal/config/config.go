package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Listen                   string `envconfig:"listen" default:""`
	Port                     string `envconfig:"port" default:"8080"`
	GoogleCloudStorageBucket string `envconfig:"google_cloud_storage_bucket"`
	MainPageSuffix           string `envconfig:"main_page_suffix" default:"index.html"`
	NotFoundPage             string `envconfig:"not_found_page"`
}

var conf Config

// LoadConf Load Configurations
func LoadConf() error {
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("config.LoadConf: failed to load conf: %w", err)
	}

	return nil
}

func Listen() string {
	return conf.Listen
}

func Port() string {
	return conf.Port
}

func GoogleCloudStorageBucket() string {
	return conf.GoogleCloudStorageBucket
}

func MainPageSuffix() string {
	return conf.MainPageSuffix
}

func NotFoundPage() string {
	return conf.NotFoundPage
}
