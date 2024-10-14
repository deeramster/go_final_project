package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port     string `envconfig:"TODO_PORT" default:"7540"`
	DBFile   string `envconfig:"TODO_DBFILE" default:"scheduler.db"`
	Password string `envconfig:"TODO_PASSWORD" required:"true"`
}

var AppConfig Config

// LoadConfig loads environment variables and maps them to the Config struct
func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on system environment variables")
	}

	// Parse environment variables into AppConfig
	err = envconfig.Process("", &AppConfig)
	if err != nil {
		log.Fatal("Error processing environment variables: ", err)
	}
}
