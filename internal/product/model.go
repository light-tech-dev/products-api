package product

import "time"

// Product يمثل منتجًا.
//
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null;index" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	SKU         string    `gorm:"size:100;uniqueIndex;not null" json:"sku"`
	Price       float64   `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	Category    string    `gorm:"size:100;index" json:"category"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
