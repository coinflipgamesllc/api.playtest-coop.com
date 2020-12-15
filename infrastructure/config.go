package infrastructure

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("dotenv")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}
}
