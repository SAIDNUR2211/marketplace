package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var AppSettings Configs

func ReadSettings() error {
	fmt.Println("Starting reading settings file")

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(".env"); err != nil {

		logger.Info().Msg("No .env file found, relying on environment variables.")
	}

	configFile, err := os.Open("internal/configs/configs.json")
	if err != nil {
		return fmt.Errorf("couldn't open config file: %w", err)
	}
	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			log.Fatal("couldn't close config file: ", err.Error())
		}
	}(configFile)

	fmt.Println("Starting decoding settings file")
	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return fmt.Errorf("couldn't decode settings json file: %w", err)
	}
	return nil
}

//release
