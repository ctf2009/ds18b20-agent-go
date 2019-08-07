package agent

import (
	"github.com/spf13/viper"
)

type Config struct {
	PORT         string
	DS18B20_ROOT string
	STORE_DIR    string
}

func NewConfig() (*Config, error) {

	var config *Config

	viper.SetConfigName("agent-config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./test")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DS18B20_ROOT", "/sys/devices/w1_bus_master1")
	viper.SetDefault("STORE_DIR", "./store")

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
