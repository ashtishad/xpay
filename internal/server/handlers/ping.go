package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	l *slog.Logger
}

func NewPingHandler(logger *slog.Logger) *PingHandler {
	return &PingHandler{
		l: logger,
	}
}

func (h *PingHandler) PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
