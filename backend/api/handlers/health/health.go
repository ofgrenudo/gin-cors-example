package health

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	/**
	Used to see if you can connect to the server from wherever you are in the world.
	*/
	slog.Debug("Ping Recieved")
	c.JSON(http.StatusOK, "Pong")
}
