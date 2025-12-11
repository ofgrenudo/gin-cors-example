package health

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ofgrenudo/gin-example/helpers"
)

func Ping(c *gin.Context) {
	/**
	Used to see if you can connect to the server from wherever you are in the world.
	*/
	slog.Debug("Ping Recieved")
	c.JSON(http.StatusOK, "Pong")
}

func AuthCheck(c *gin.Context) {
	user, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"auth":     true,
		"username": user.Username,
	})
}
