package util

import (
	"github.com/spf13/viper"
)
//Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver 	  string `mapstructure:"DB_DRIVER"`
	DBSource 	  string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

//LoadConfig reads configuration from a file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can be json,xml, etc

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
