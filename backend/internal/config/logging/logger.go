package logging

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func convertStringToSlogLevel(logLevelToBe string) slog.Level {
	logLevel := slog.LevelInfo

	switch strings.ToUpper(logLevelToBe) {
	case "INFO":
		logLevel = slog.LevelInfo
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		// Default case handles any value that doesn't match above
		log.Printf("Error: Unable to interpret the Log Level '%v'. Please confirm the log level in the .env is set to INFO, DEBUG, WARN, or ERROR.", logLevelToBe)
		log.Printf("Defaulting to INFO")
		logLevel = slog.LevelInfo
	}

	return logLevel
}

func InitializeGlobalLogger(logFilePath string, minimumLevel string) func() {
	/**
	 * Initialize Global Logger, will initialize a global slogger to write to a specific file,
	 * in a specific way. It also makes sure that the log file is closed when the application exits
	 */
	// this will open the file, with open the file as write only, but then optionally
	// create, or append if needed. It will also expect / create the file with the 0666 permission.
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open the log file at: %v", err)
	}

	slogLevel := convertStringToSlogLevel(minimumLevel)

	// there are other options, like a traditional log handler, but i think that a JSON handler is fine
	// for what we are trying to do.
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: (slogLevel == slog.LevelDebug),
	})

	slog.SetDefault(slog.New(handler))

	return func() {
		if file != nil {
			slog.Info("Closing the log file...")
			_ = file.Close()
		}
	}
}
