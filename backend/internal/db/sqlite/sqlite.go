package sqlite

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/ofgrenudo/gin-example/internal/config/env"
	"github.com/ofgrenudo/gin-example/internal/db/models/users"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log/slog"
)

var GlobalDB *gorm.DB

func InitDB(dbPath string) {
	var err error

	dbDir := getDirFromPath(dbPath)
	if dbDir != "" {
		if err := os.MkdirAll(dbDir, 0o755); err != nil {
			slog.Error("Failed to create database directory", slog.Any("error", err), slog.String("path", dbDir))
			os.Exit(1)
		}
	}

	GlobalDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Database connection established and ready!")

	// Auto-migrate User model (Username + APIKey)
	if err := GlobalDB.AutoMigrate(&users.User{}); err != nil {
		slog.Error("Failed to auto-migrate database schema", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Database migrations completed successfully.")

	createDefaultUser()
	// todo(jwb): wireup auth middleware (this is done in main(), not here)
}

// Helper to extract directory from path
func getDirFromPath(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			// Found a file extension, so we assume it's a file path
			// Look for the last slash
			lastSlash := -1
			for j := len(path) - 1; j >= 0; j-- {
				if path[j] == '/' {
					lastSlash = j
					break
				}
			}
			if lastSlash != -1 {
				return path[:lastSlash]
			}
			return ""
		}
	}
	return "" // Assume root if no slashes are present
}

func createDefaultUser() {
	/**
	createDefaultUser will go through and create a default user, according to the environment variables configured.
	*/

	defaultUserName := env.GlobalConfig.BackendDefaultUserName
	defaultAPIKey := env.GlobalConfig.BackendDefaultAPIKey // you define this in your env config

	if defaultUserName == "" {
		slog.Info("No default username configured, skipping default user creation")
		return
	}

	// Check if user already exists
	var count int64
	if err := GlobalDB.Model(&users.User{}).
		Where("username = ?", defaultUserName).
		Count(&count).Error; err != nil {
		slog.Error("Failed checking for default user", slog.Any("error", err))
		os.Exit(1)
	}

	if count > 0 {
		slog.Info("Default user already exists", slog.String("username", defaultUserName))
		return
	}

	// If no API key is configured in env, generate one
	if defaultAPIKey == "" {
		var err error
		defaultAPIKey, err = generateAPIKey()
		if err != nil {
			slog.Error("Failed to generate API key for default user", slog.Any("error", err))
			os.Exit(1)
		}
		// Only log a prefix so you can correlate without leaking the key
		keyPrefix := defaultAPIKey
		if len(keyPrefix) > 4 {
			keyPrefix = keyPrefix[:4]
		}
		slog.Info("Generated API key for default user",
			slog.String("username", defaultUserName),
			slog.String("api_key_prefix", keyPrefix),
		)
	} else {
		// If you insist on env-provided key, still only log prefix
		keyPrefix := defaultAPIKey
		if len(keyPrefix) > 4 {
			keyPrefix = keyPrefix[:4]
		}
		slog.Info("Using configured API key for default user",
			slog.String("username", defaultUserName),
			slog.String("api_key_prefix", keyPrefix),
		)
	}

	defaultUser := users.User{
		Username: defaultUserName,
		APIKey:   defaultAPIKey,
	}

	if err := GlobalDB.Create(&defaultUser).Error; err != nil {
		slog.Error("Failed to create default user", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Default user created successfully", slog.String("username", defaultUserName))
}

func generateAPIKey() (string, error) {
	// 32 bytes = 256-bit key
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
