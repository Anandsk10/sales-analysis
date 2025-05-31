package services

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test helper functions
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func createTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Suppress logs during tests
	return logger
}

// AnyTime is a matcher for time.Time values in SQL mocks
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// AnalyticsService Tests
func TestAnalyticsService_GetTotalRevenue(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		expectedQuery := `SELECT 
            COALESCE\(SUM\(oi\.quantity_sold \* oi\.unit_price \* \(1 - oi\.discount\)\), 0\) as revenue,
            COUNT\(DISTINCT o\.order_id\) as count
        FROM orders o
        JOIN order_items oi ON o\.order_id = oi\.order_id
        WHERE o\.date_of_sale BETWEEN \$1 AND \$2`

		rows := sqlmock.NewRows([]string{"revenue", "count"}).
			AddRow(10000.50, 25)

		mock.ExpectQuery(expectedQuery).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		result, err := service.GetTotalRevenue(startDate, endDate)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10000.50, result.Revenue)
		assert.Equal(t, int64(25), result.Count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestAnalyticsService_GetRevenueByProduct(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"product_id", "product_name", "revenue", "count"}).
			AddRow("P001", "Product 1", 5000.25, 10).
			AddRow("P002", "Product 2", 3000.75, 8)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		results, err := service.GetRevenueByProduct(startDate, endDate)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "P001", results[0].ProductID)
		assert.Equal(t, "Product 1", results[0].ProductName)
		assert.Equal(t, 5000.25, results[0].Revenue)
		assert.Equal(t, int64(10), results[0].Count)
	})
}

func TestAnalyticsService_GetRevenueByCategory(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"category", "revenue", "count"}).
			AddRow("Electronics", 15000.00, 20).
			AddRow("Clothing", 8000.50, 15)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		results, err := service.GetRevenueByCategory(startDate, endDate)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Electronics", results[0].Category)
		assert.Equal(t, 15000.00, results[0].Revenue)
		assert.Equal(t, int64(20), results[0].Count)
	})
}

func TestAnalyticsService_GetTopProducts(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	limit := 5

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"product_id", "product_name", "category", "total_sold", "revenue"}).
			AddRow("P001", "Product 1", "Electronics", 100, 5000.00).
			AddRow("P002", "Product 2", "Clothing", 80, 3200.00)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(startDate, endDate, limit).
			WillReturnRows(rows)

		results, err := service.GetTopProducts(startDate, endDate, limit)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "P001", results[0].ProductID)
		assert.Equal(t, int64(100), results[0].TotalSold)
		assert.Equal(t, 5000.00, results[0].Revenue)
	})
}

func TestAnalyticsService_GetCustomerCount(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(150)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(DISTINCT o.customer_id)")).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		count, err := service.GetCustomerCount(startDate, endDate)

		assert.NoError(t, err)
		assert.Equal(t, int64(150), count)
	})
}

func TestAnalyticsService_GetOrderCount(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(200)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "orders"`)).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		count, err := service.GetOrderCount(startDate, endDate)

		assert.NoError(t, err)
		assert.Equal(t, int64(200), count)
	})
}

func TestAnalyticsService_GetAverageOrderValue(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"avg_value"}).AddRow(125.75)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(AVG(order_total), 0)")).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		avgValue, err := service.GetAverageOrderValue(startDate, endDate)

		assert.NoError(t, err)
		assert.Equal(t, 125.75, avgValue)
	})
}

func TestRefreshService_clearExistingData(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	csvLoader := NewCSVLoader(db, logger)
	service := NewRefreshService(db, csvLoader, logger)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM order_items").WillReturnResult(sqlmock.NewResult(0, 5))
		mock.ExpectExec("DELETE FROM orders").WillReturnResult(sqlmock.NewResult(0, 3))
		mock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectExec("DELETE FROM customers").WillReturnResult(sqlmock.NewResult(0, 2))

		err := service.clearExistingData()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
func TestRefreshService_GetRefreshStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	logger := createTestLogger()
	csvLoader := NewCSVLoader(db, logger)
	service := NewRefreshService(db, csvLoader, logger)

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "status", "start_time", "end_time", "records_count", "error_message", "created_at"}).
			AddRow(1, "success", now, now, 100, "", now).
			AddRow(2, "failed", now, now, 0, "error occurred", now)

		// Fix: Use regex to match the parameterized query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_logs" ORDER BY created_at DESC LIMIT $1`)).
			WithArgs(10).
			WillReturnRows(rows)

		logs, err := service.GetRefreshStatus()

		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		assert.Equal(t, "success", logs[0].Status)
		assert.Equal(t, "failed", logs[1].Status)
	})

}

// Benchmark Tests
func BenchmarkAnalyticsService_GetTotalRevenue(b *testing.B) {
	db, mock := setupMockDB(&testing.T{})
	logger := createTestLogger()
	service := NewAnalyticsService(db, logger)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"revenue", "count"}).AddRow(10000.50, 25)

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		service.GetTotalRevenue(startDate, endDate)
	}
}

// Integration Test Setup Example
func TestAnalyticsServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would be used with a real database connection
	// for integration testing
	t.Skip("Integration tests require real database setup")
}
