package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RabbitSource         string        `mapstructure:"RABBIT_SOURCE"`
	AllowedOrigins       string        `mapstructure:"ALLOWED_ORIGINS"`
	MediaBasePath        string        `mapstructure:"MEDIA_BASE_PATH"`
	ImagePath            string        `mapstructure:"IMAGE_PATH"`
	VideoPath            string        `mapstructure:"VIDEO_PATH"`
}

// LoadConfig reads configuration from file or envirnment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("RABBIT_SOURCE")
	viper.BindEnv("ALLOWED_ORIGINS")
	viper.BindEnv("MEDIA_BASE_PATH")
	viper.BindEnv("IMAGE_PATH")
	viper.BindEnv("VIDEO_PATH")

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
