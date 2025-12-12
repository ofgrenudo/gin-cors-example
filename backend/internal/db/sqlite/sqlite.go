package sqlite

import (
	"log"
	"os"

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
			log.Fatal("Failed to auto-migrate database schema", err)
		}
	}

	GlobalDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("error", err))
		log.Fatal("Failed to connect to database", err)
	}

	slog.Info("Database connection established and ready!")

	// Auto-migrate User model (Username + APIKey)
	if err := GlobalDB.AutoMigrate(&users.User{}); err != nil {
		slog.Error("Failed to auto-migrate database schema", slog.Any("error", err))
		log.Fatal("Failed to auto-migrate database schema", err)
	}
	slog.Info("Database migrations completed successfully.")
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
