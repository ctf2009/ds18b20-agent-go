package config

type Config struct {
	PORT  string
}

//TODO: Add Viper
func New() (*Config, error)  {

	config := &Config{
		PORT: "8080",
	}

	return config, nil
}