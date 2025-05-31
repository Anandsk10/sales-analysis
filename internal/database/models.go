package database

import (
	"time"
)

type Customer struct {
	ID        string    `gorm:"primaryKey;column:customer_id" json:"customer_id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Orders []Order `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

type Product struct {
	ID        string    `gorm:"primaryKey;column:product_id" json:"product_id"`
	Name      string    `gorm:"not null" json:"name"`
	Category  string    `gorm:"not null;index" json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	OrderItems []OrderItem `gorm:"foreignKey:ProductID" json:"order_items,omitempty"`
}

type Order struct {
	ID            string    `gorm:"primaryKey;column:order_id" json:"order_id"`
	CustomerID    string    `gorm:"not null;index" json:"customer_id"`
	Region        string    `gorm:"not null;index" json:"region"`
	DateOfSale    time.Time `gorm:"not null;index" json:"date_of_sale"`
	PaymentMethod string    `json:"payment_method"`
	ShippingCost  float64   `gorm:"type:decimal(10,2)" json:"shipping_cost"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Customer   Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      string    `gorm:"not null;index" json:"order_id"`
	ProductID    string    `gorm:"not null;index" json:"product_id"`
	QuantitySold int       `gorm:"not null" json:"quantity_sold"`
	UnitPrice    float64   `gorm:"not null;type:decimal(10,2)" json:"unit_price"`
	Discount     float64   `gorm:"type:decimal(5,4)" json:"discount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Order   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type RefreshLog struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Status       string     `gorm:"not null" json:"status"` // success, failed, in_progress
	StartTime    time.Time  `gorm:"not null" json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	RecordsCount int        `json:"records_count"`
	ErrorMessage string     `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
}
