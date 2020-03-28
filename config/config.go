package config

import (
	"github.com/spf13/viper"
	"github.com/creasty/defaults"
	"log"
)

type Configuration struct {
	Server ServerConfiguration
}

type ServerConfiguration struct {
	Port int `default:"8000"`
}

var (
	Config Configuration
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func (s *ServerConfiguration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)

	type plain ServerConfiguration
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}