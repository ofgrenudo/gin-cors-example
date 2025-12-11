package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	/**
	 * Config should all configurations relevant to the application. If you need
	 * something configured or to read something from the ENV, please put it here.
	 */
	LogPath     string
	LogLevel    string
	BackendPort int
}

// Global config is a global variable to hold the initialized configuration struct.
var GlobalConfig *Config

func InitializeConfig() error {
	if err := godotenv.Load(); err != nil {
		// this must be log, because we use the env to configure slog.
		log.Println("No .env file found, continuing with system environment variables")
	}

	cfg := Config{
		LogPath:     os.Getenv("LOG_PATH"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		BackendPort: stringEnvToInt("BACKEND_PORT"),
	}

	// pass a reference to our global config.
	GlobalConfig = &cfg
	return nil
}

func stringEnvToInt(envName string) int {
	pendingEnv, err := strconv.Atoi(os.Getenv(envName))
	if err != nil {
		log.Println(err)
		log.Fatalf("Unable to get or convert %v from .env", envName)
	}

	return pendingEnv
}
