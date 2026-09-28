package item

import (
	"errors"
	"strings"
	"time"
)

type CreateItemRequest struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Notes     string  `json:"notes"`
}

func (r *CreateItemRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if r.Quantity < 0 {
		return errors.New("quantity must be >= 0")
	}
	if r.UnitPrice < 0 {
		return errors.New("unit_price must be >= 0")
	}
	return nil
}

type ItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Notes     string  `json:"notes"`
	CreatedAt string  `json:"created_at"`
}

func FromItem(i *Item) *ItemResponse {
	return &ItemResponse{
		ID:        i.ID,
		ProductID: i.ProductID,
		Name:      i.Name,
		Quantity:  i.Quantity,
		UnitPrice: i.UnitPrice,
		Notes:     i.Notes,
		CreatedAt: i.CreatedAt.Format(time.RFC3339),
	}
}
