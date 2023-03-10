package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	JwtSecret     string        `mapstructure:"JWT_SECRET"`
	JwtDuration   time.Duration `mapstructure:"JWT_DURATION"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)
	return
}
