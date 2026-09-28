package item

import "time"

type Item struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Quantity  int       `gorm:"default:1" json:"quantity"`
	UnitPrice float64   `gorm:"type:decimal(10,2);default:0" json:"unit_price"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Item) TableName() string {
	return "items"
}
