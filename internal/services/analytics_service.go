package services

import (
	"sales-analysis-system/internal/database"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AnalyticsService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

type RevenueResult struct {
	Revenue float64 `json:"revenue"`
	Count   int64   `json:"count"`
}

type ProductRevenueResult struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Revenue     float64 `json:"revenue"`
	Count       int64   `json:"count"`
}

type CategoryRevenueResult struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Count    int64   `json:"count"`
}

type RegionRevenueResult struct {
	Region  string  `json:"region"`
	Revenue float64 `json:"revenue"`
	Count   int64   `json:"count"`
}

type TopProductResult struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	TotalSold   int64   `json:"total_sold"`
	Revenue     float64 `json:"revenue"`
}

func NewAnalyticsService(db *gorm.DB, logger *logrus.Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger,
	}
}

func (a *AnalyticsService) GetTotalRevenue(startDate, endDate time.Time) (*RevenueResult, error) {
	var result RevenueResult

	query := `
        SELECT 
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue,
            COUNT(DISTINCT o.order_id) as count
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        WHERE o.date_of_sale BETWEEN ? AND ?
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&result).Error
	return &result, err
}

func (a *AnalyticsService) GetRevenueByProduct(startDate, endDate time.Time) ([]ProductRevenueResult, error) {
	var results []ProductRevenueResult

	query := `
        SELECT 
            p.product_id,
            p.name as product_name,
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue,
            COUNT(DISTINCT o.order_id) as count
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        JOIN products p ON oi.product_id = p.product_id
        WHERE o.date_of_sale BETWEEN ? AND ?
        GROUP BY p.product_id, p.name
        ORDER BY revenue DESC
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

func (a *AnalyticsService) GetRevenueByCategory(startDate, endDate time.Time) ([]CategoryRevenueResult, error) {
	var results []CategoryRevenueResult

	query := `
        SELECT 
            p.category,
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue,
            COUNT(DISTINCT o.order_id) as count
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        JOIN products p ON oi.product_id = p.product_id
        WHERE o.date_of_sale BETWEEN ? AND ?
        GROUP BY p.category
        ORDER BY revenue DESC
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

func (a *AnalyticsService) GetRevenueByRegion(startDate, endDate time.Time) ([]RegionRevenueResult, error) {
	var results []RegionRevenueResult

	query := `
        SELECT 
            o.region,
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue,
            COUNT(DISTINCT o.order_id) as count
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        WHERE o.date_of_sale BETWEEN ? AND ?
        GROUP BY o.region
        ORDER BY revenue DESC
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&results).Error
	return results, err
}

func (a *AnalyticsService) GetTopProducts(startDate, endDate time.Time, limit int) ([]TopProductResult, error) {
	var results []TopProductResult

	query := `
        SELECT 
            p.product_id,
            p.name as product_name,
            p.category,
            COALESCE(SUM(oi.quantity_sold), 0) as total_sold,
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        JOIN products p ON oi.product_id = p.product_id
        WHERE o.date_of_sale BETWEEN ? AND ?
        GROUP BY p.product_id, p.name, p.category
        ORDER BY total_sold DESC
        LIMIT ?
    `

	err := a.db.Raw(query, startDate, endDate, limit).Scan(&results).Error
	return results, err
}

func (a *AnalyticsService) GetTopProductsByCategory(startDate, endDate time.Time, category string, limit int) ([]TopProductResult, error) {
	var results []TopProductResult

	query := `
        SELECT 
            p.product_id,
            p.name as product_name,
            p.category,
            COALESCE(SUM(oi.quantity_sold), 0) as total_sold,
            COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) as revenue
        FROM orders o
        JOIN order_items oi ON o.order_id = oi.order_id
        JOIN products p ON oi.product_id = p.product_id
        WHERE o.date_of_sale BETWEEN ? AND ? AND p.category = ?
        GROUP BY p.product_id, p.name, p.category
        ORDER BY total_sold DESC
        LIMIT ?
    `

	err := a.db.Raw(query, startDate, endDate, category, limit).Scan(&results).Error
	return results, err
}

func (a *AnalyticsService) GetCustomerCount(startDate, endDate time.Time) (int64, error) {
	var count int64

	query := `
        SELECT COUNT(DISTINCT o.customer_id)
        FROM orders o
        WHERE o.date_of_sale BETWEEN ? AND ?
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&count).Error
	return count, err
}

func (a *AnalyticsService) GetOrderCount(startDate, endDate time.Time) (int64, error) {
	var count int64

	err := a.db.Model(&database.Order{}).
		Where("date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error

	return count, err
}

func (a *AnalyticsService) GetAverageOrderValue(startDate, endDate time.Time) (float64, error) {
	var avgValue float64

	query := `
        SELECT COALESCE(AVG(order_total), 0) as avg_value
        FROM (
            SELECT 
                o.order_id,
                SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)) as order_total
            FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            WHERE o.date_of_sale BETWEEN ? AND ?
            GROUP BY o.order_id
        ) as order_totals
    `

	err := a.db.Raw(query, startDate, endDate).Scan(&avgValue).Error
	return avgValue, err
}
