package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type XConfig struct {
	Server struct {
		Port string `mapstructure:"SERVER_PORT"`
		Host string `mapstructure:"SERVER_HOST"`
	}
	Database struct {
		DBPassword string `mapstructure:"DB_PASSWORD"`
	}
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}
type Config struct {
	Env    string `mapstructure:"APP_ENV"`
	DBName string `mapstructure:"DB_NAME"`
	DBUser string `mapstructure:"DB_USERNAME"`
	DBPass string `mapstructure:"DB_PASSWORD"`
	DBAddr string `mapstructure:"DB_ADDR"`
	Port   string `mapstructure:"PORT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
