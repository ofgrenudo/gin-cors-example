package sqlite

import (
	"github.com/ofgrenudo/gin-example/internal/db/models/users"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log/slog"
	"os"
)

var GlobalDB *gorm.DB

func InitDB(dbPath string) {
	var err error

	dbDir := getDirFromPath(dbPath)
	if dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
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

	err = GlobalDB.AutoMigrate(&users.User{})
	if err != nil {
		slog.Error("Failed to auto-migrate database schema", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Database migrations completed successfully.")

	//todo(jwb): add a default user.
	//todo(jwb): wireup auth middleware
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
	return "" // Assume root directory if no slashes are present
}
