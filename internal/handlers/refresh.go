package handlers

import (
	"net/http"
	"sales-analysis-system/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RefreshHandler struct {
	service *services.RefreshService
	logger  *logrus.Logger
}

func NewRefreshHandler(service *services.RefreshService, logger *logrus.Logger) *RefreshHandler {
	return &RefreshHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RefreshHandler) TriggerRefresh(c *gin.Context) {
	filePath := c.DefaultQuery("file_path", "data/sales_data.csv")

	// Run refresh in background
	go func() {
		if err := h.service.RefreshData(filePath); err != nil {
			h.logger.Error("Background refresh failed: ", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Data refresh triggered successfully",
		"status":    "in_progress",
		"file_path": filePath,
	})
}

func (h *RefreshHandler) GetRefreshStatus(c *gin.Context) {
	logs, err := h.service.GetRefreshStatus()
	if err != nil {
		h.logger.Error("Failed to get refresh status: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get refresh status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}
