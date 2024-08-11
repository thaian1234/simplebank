package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

type Environment interface {
	GetConfig() (config Config, err error)
}

type ProductionEnv struct{}
type DevEnv struct {
	path string
}

func NewEnvironment(path string) Environment {
	isProductionEnv := os.Getenv("APP_ENV") == "production"

	if isProductionEnv {
		return &ProductionEnv{}
	}
	return &DevEnv{
		path: path,
	}
}

func (pe ProductionEnv) GetConfig() (config Config, err error) {
	return Config{
		DBDriver:            os.Getenv("DB_DRIVER"),
		DBSource:            os.Getenv("DB_SOURCE"),
		ServerAddress:       os.Getenv("SERVER_ADDRESS"),
		TokenSymmetricKey:   os.Getenv("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration: time.Duration(15) * time.Minute,
	}, nil
}

func (de DevEnv) GetConfig() (config Config, err error) {
	viper.AddConfigPath(de.path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //json, xml
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
