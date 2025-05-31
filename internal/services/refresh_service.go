package services

import (
	"sales-analysis-system/internal/database"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RefreshService struct {
	db        *gorm.DB
	csvLoader *CSVLoader
	logger    *logrus.Logger
}

func NewRefreshService(db *gorm.DB, csvLoader *CSVLoader, logger *logrus.Logger) *RefreshService {
	return &RefreshService{
		db:        db,
		csvLoader: csvLoader,
		logger:    logger,
	}
}

func (r *RefreshService) RefreshData(filePath string) error {
	// Log refresh start
	refreshLog := database.RefreshLog{
		Status:    "in_progress",
		StartTime: time.Now(),
	}
	r.db.Create(&refreshLog)

	r.logger.Info("Starting data refresh from: ", filePath)

	// Clear existing data
	if err := r.clearExistingData(); err != nil {
		r.updateRefreshLog(refreshLog.ID, "failed", 0, err.Error())
		return err
	}

	// Load new data
	if err := r.csvLoader.LoadFromCSV(filePath); err != nil {
		r.updateRefreshLog(refreshLog.ID, "failed", 0, err.Error())
		return err
	}

	// Count records
	var count int64
	r.db.Model(&database.Order{}).Count(&count)

	// Update refresh log
	r.updateRefreshLog(refreshLog.ID, "success", int(count), "")

	r.logger.Info("Data refresh completed successfully")
	return nil
}

func (r *RefreshService) clearExistingData() error {
	r.logger.Info("Clearing existing data...")

	if err := r.db.Exec("DELETE FROM order_items").Error; err != nil {
		return err
	}
	if err := r.db.Exec("DELETE FROM orders").Error; err != nil {
		return err
	}
	if err := r.db.Exec("DELETE FROM products").Error; err != nil {
		return err
	}
	if err := r.db.Exec("DELETE FROM customers").Error; err != nil {
		return err
	}

	return nil
}

func (r *RefreshService) updateRefreshLog(id uint, status string, recordsCount int, errorMessage string) {
	endTime := time.Now()
	r.db.Model(&database.RefreshLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"end_time":      &endTime,
		"records_count": recordsCount,
		"error_message": errorMessage,
	})
}

func (r *RefreshService) GetRefreshStatus() ([]database.RefreshLog, error) {
	var logs []database.RefreshLog
	err := r.db.Order("created_at DESC").Limit(10).Find(&logs).Error
	return logs, err
}
