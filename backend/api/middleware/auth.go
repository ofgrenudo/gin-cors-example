package middleware

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
	// UserKey is the key used to store the authenticated user in Gin's context.
	UserKey = "authenticatedUser"
)

// AuthMiddleware creates a middleware that checks for a valid API Key.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract API Key from the Authorization header (Bearer or API-Key schema)
		authHeader := c.GetHeader("Authorization")

		apiKey := ""
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)

			// Check for 'Bearer' or 'API-Key' scheme
			if len(parts) == 2 && (parts[0] == "Bearer" || parts[0] == "API-Key") {
				apiKey = parts[1]
			}
		}

		if apiKey == "" {
			slog.Warn("Authentication failed: Missing API Key in request header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key is required"})
			return
		}

		// 2. Validate API Key against the database using the global connection
		var user users.User

		// Use sqlite.GlobalDB to perform the query
		result := sqlite.GlobalDB.Where("api_key = ?", apiKey).First(&user)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Log the prefix, not the full key, for security
				keyPrefix := "N/A"
				if len(apiKey) > 4 {
					keyPrefix = apiKey[:4]
				}
				slog.Warn("Authentication failed: Invalid API Key", slog.String("key_prefix", keyPrefix))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
				return
			}

			// Handle other database errors
			slog.Error("Database error during authentication", slog.Any("error", result.Error))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error during authentication"})
			return
		}

		// 3. Success: Store the user object in the context
		slog.Debug("Authentication successful", slog.String("username", user.Username))
		c.Set(UserKey, &user)

		// 4. Continue to the next handler
		c.Next()
	}
}
