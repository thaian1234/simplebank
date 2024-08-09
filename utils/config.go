package utils

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

type Environment interface {
	GetConfig() (config Config, err error)
}

type ProductionEnv struct{}
type DevEnv struct {
	path string
}

func NewEnviroment(path string) Environment {
	isProductionEnv := os.Getenv("APP_ENV") == "production"

	if isProductionEnv {
		return &ProductionEnv{}
	}
	return &DevEnv{
		path: path,
	}
}

func (pe ProductionEnv) GetConfig() (cofnig Config, err error) {
	return Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBSource:      os.Getenv("DB_SOURCE"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
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
