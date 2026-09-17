package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheck godoc
// @Summary      Health check
// @Tags         monitoring
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h Handler) HealthCheck(c *gin.Context) {
	db, err := h.DB.DB()

	if err != nil {
		h.Logger.Error("database connection failed", zap.Error(err))
		ErrorResponse(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	err = db.Ping()

	if err != nil {
		h.Logger.Error("database connection failed", zap.Error(err))
		ErrorResponse(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	h.Logger.Info("health check passed")
	SuccessResponse(c, http.StatusOK, gin.H{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": time.Now(),
	})
}
