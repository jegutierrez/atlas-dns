package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller holds ping request handler.
type Controller interface {
	PingHandler(ctx *gin.Context)
}

// NewPongController builds ping's controller.
func NewPongController() Controller {
	return &PongController{}
}

// PongController holds the ping http request handler.
type PongController struct{}

// PingHandler responds to ping requests.
// Is intended to be used to check if the application is alive.
func (p *PongController) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
