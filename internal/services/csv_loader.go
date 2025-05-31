package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sales-analysis-system/internal/database"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CSVLoader struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewCSVLoader(db *gorm.DB, logger *logrus.Logger) *CSVLoader {
	return &CSVLoader{
		db:     db,
		logger: logger,
	}
}

func (c *CSVLoader) LoadFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	c.logger.Info("CSV Headers: ", headers)

	// Start transaction
	tx := c.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	recordCount := 0
	batchSize := 1000
	var customers []database.Customer
	var products []database.Product
	var orders []database.Order
	var orderItems []database.OrderItem

	customerMap := make(map[string]database.Customer)
	productMap := make(map[string]database.Product)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Parse CSV record
		orderID := record[0]
		productID := record[1]
		customerID := record[2]
		productName := record[3]
		category := record[4]
		region := record[5]
		dateStr := record[6]
		quantityStr := record[7]
		unitPriceStr := record[8]
		discountStr := record[9]
		shippingCostStr := record[10]
		paymentMethod := record[11]
		customerName := record[12]
		customerEmail := record[13]
		customerAddress := record[14]

		// Parse numeric values
		quantity, _ := strconv.Atoi(quantityStr)
		unitPrice, _ := strconv.ParseFloat(unitPriceStr, 64)
		discount, _ := strconv.ParseFloat(discountStr, 64)
		shippingCost, _ := strconv.ParseFloat(shippingCostStr, 64)

		// Parse date
		dateOfSale, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.logger.Warn("Failed to parse date: ", dateStr)
			dateOfSale = time.Now()
		}

		// Create customer if not exists
		if _, exists := customerMap[customerID]; !exists {
			customer := database.Customer{
				ID:      customerID,
				Name:    customerName,
				Email:   customerEmail,
				Address: customerAddress,
			}
			customerMap[customerID] = customer
			customers = append(customers, customer)
		}

		// Create product if not exists
		if _, exists := productMap[productID]; !exists {
			product := database.Product{
				ID:       productID,
				Name:     productName,
				Category: category,
			}
			productMap[productID] = product
			products = append(products, product)
		}

		// Create order
		order := database.Order{
			ID:            orderID,
			CustomerID:    customerID,
			Region:        region,
			DateOfSale:    dateOfSale,
			PaymentMethod: paymentMethod,
			ShippingCost:  shippingCost,
		}
		orders = append(orders, order)

		// Create order item
		orderItem := database.OrderItem{
			OrderID:      orderID,
			ProductID:    productID,
			QuantitySold: quantity,
			UnitPrice:    unitPrice,
			Discount:     discount,
		}
		orderItems = append(orderItems, orderItem)

		recordCount++

		// Batch insert
		if recordCount%batchSize == 0 {
			if err := c.batchInsert(tx, customers, products, orders, orderItems); err != nil {
				tx.Rollback()
				return err
			}
			customers = customers[:0]
			products = products[:0]
			orders = orders[:0]
			orderItems = orderItems[:0]
		}
	}

	// Insert remaining records
	if len(customers) > 0 || len(products) > 0 || len(orders) > 0 || len(orderItems) > 0 {
		if err := c.batchInsert(tx, customers, products, orders, orderItems); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	c.logger.Info(fmt.Sprintf("Successfully loaded %d records from CSV", recordCount))
	return nil
}

func (c *CSVLoader) batchInsert(tx *gorm.DB, customers []database.Customer, products []database.Product, orders []database.Order, orderItems []database.OrderItem) error {
	if len(customers) > 0 {
		if err := tx.CreateInBatches(customers, 100).Error; err != nil {
			c.logger.Error("Failed to insert customers: ", err)
		}
	}

	if len(products) > 0 {
		if err := tx.CreateInBatches(products, 100).Error; err != nil {
			c.logger.Error("Failed to insert products: ", err)
		}
	}

	if len(orders) > 0 {
		if err := tx.CreateInBatches(orders, 100).Error; err != nil {
			return fmt.Errorf("failed to insert orders: %w", err)
		}
	}

	if len(orderItems) > 0 {
		if err := tx.CreateInBatches(orderItems, 100).Error; err != nil {
			return fmt.Errorf("failed to insert order items: %w", err)
		}
	}

	return nil
}
