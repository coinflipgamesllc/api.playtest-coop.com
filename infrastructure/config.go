package infrastructure

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()

	// Only pull from .env in development
	environment := os.Getenv("ENVIRONMENT")
	if environment == "development" || environment == "" {
		viper.SetConfigType("dotenv")
		viper.SetConfigName(".env")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Failed to parse config file: %v\n", err)
		}
	}
}
