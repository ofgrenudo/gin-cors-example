package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	// Correctly import the package containing the GlobalDB instance
	"github.com/ofgrenudo/gin-example/internal/db/sqlite"

	// Correctly import the package containing the User model
	"github.com/ofgrenudo/gin-example/internal/db/models/users"

	"gorm.io/gorm"
)

const (
	UserKey = "authenticatedUser"
)

func AuthMiddleware() gin.HandlerFunc {
	/**
	Checks for a valid API Key.
	*/
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		apiKey := ""
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && (parts[0] == "Bearer" || parts[0] == "API-Key") {
				apiKey = parts[1]
			}
		}

		if apiKey == "" {
			slog.Warn("Authentication failed: Missing API Key in request header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You have not authenticated with this service yet."})
			return
		}

		// validate API Key against the database using the global connection
		var user users.User
		result := sqlite.GlobalDB.Where("api_key = ?", apiKey).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				keyPrefix := "N/A"
				if len(apiKey) > 4 {
					keyPrefix = apiKey[:4]
				}
				slog.Warn("Authentication failed: Invalid API Key", slog.String("key_prefix", keyPrefix))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
				return
			}

			slog.Error("Database error during authentication", slog.Any("error", result.Error))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error during authentication"})
			return
		}

		slog.Debug("Authentication successful", slog.String("username", user.Username))
		c.Set(UserKey, &user)
		c.Next()
	}
}
