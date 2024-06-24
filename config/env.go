package config

type Config struct {
	Port       int
	MaxClients int
}

var Envs = initConfig()

// Initialize the configuration with default values
func initConfig() *Config {
	return &Config{
		Port:       9000,
		MaxClients: 10,
	}
}
