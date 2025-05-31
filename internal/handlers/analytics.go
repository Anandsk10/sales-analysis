package handlers

import (
	"net/http"
	"sales-analysis-system/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
	logger  *logrus.Logger
}

func NewAnalyticsHandler(service *services.AnalyticsService, logger *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AnalyticsHandler) parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startDateStr := c.DefaultQuery("start_date", "2023-01-01")
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}

func (h *AnalyticsHandler) GetTotalRevenue(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	result, err := h.service.GetTotalRevenue(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get total revenue: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total revenue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetRevenueByProduct(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	results, err := h.service.GetRevenueByProduct(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get revenue by product: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate revenue by product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetRevenueByCategory(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	results, err := h.service.GetRevenueByCategory(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get revenue by category: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate revenue by category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetRevenueByRegion(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	results, err := h.service.GetRevenueByRegion(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get revenue by region: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate revenue by region"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetRevenueTrends(c *gin.Context) {
	// This would require more complex aggregation by month/quarter/year
	// For now, return a simple message
	c.JSON(http.StatusOK, gin.H{
		"message":    "Revenue trends endpoint - implementation depends on specific requirements",
		"suggestion": "Use other endpoints with different date ranges to analyze trends",
	})
}

func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	results, err := h.service.GetTopProducts(startDate, endDate, limit)
	if err != nil {
		h.logger.Error("Failed to get top products: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"limit": limit,
	})
}

func (h *AnalyticsHandler) GetTopProductsByCategory(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category parameter is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	results, err := h.service.GetTopProductsByCategory(startDate, endDate, category, limit)
	if err != nil {
		h.logger.Error("Failed to get top products by category: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top products by category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     results,
		"category": category,
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"limit": limit,
	})
}

func (h *AnalyticsHandler) GetTopProductsByRegion(c *gin.Context) {
	// Similar implementation to GetTopProductsByCategory but for region
	c.JSON(http.StatusOK, gin.H{
		"message": "Top products by region endpoint - similar to category implementation",
	})
}

func (h *AnalyticsHandler) GetCustomerCount(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	count, err := h.service.GetCustomerCount(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get customer count: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get customer count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"customer_count": count,
		},
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetOrderCount(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	count, err := h.service.GetOrderCount(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get order count: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"order_count": count,
		},
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}

func (h *AnalyticsHandler) GetAverageOrderValue(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	avgValue, err := h.service.GetAverageOrderValue(startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get average order value: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get average order value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"average_order_value": avgValue,
		},
		"date_range": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}
