package config

import (
	"fmt"
	"log"
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
	AppEnv string `mapstructure:"APP_ENV"`
	DB     struct {
		Name string `mapstructure:"DB_NAME"`
		User string `mapstructure:"DB_USERNAME"`
		Pass string `mapstructure:"DB_PASSWORD"`
		Addr string `mapstructure:"DB_ADDR"`
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("error reading config file: %w \n", err))
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return
}
