package product

import (
	"errors"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════
// CreateProductRequest
// ═══════════════════════════════════════════════

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

func (r *CreateProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if len(r.Name) > 255 {
		return errors.New("name too long (max 255)")
	}
	if strings.TrimSpace(r.SKU) == "" {
		return errors.New("sku is required")
	}
	if r.Price < 0 {
		return errors.New("price must be >= 0")
	}
	if r.Stock < 0 {
		return errors.New("stock must be >= 0")
	}
	return nil
}

func (r *CreateProductRequest) ToModel() *Product {
	return &Product{
		Name:        strings.TrimSpace(r.Name),
		Description: strings.TrimSpace(r.Description),
		SKU:         strings.TrimSpace(strings.ToUpper(r.SKU)),
		Price:       r.Price,
		Stock:       r.Stock,
		Category:    strings.TrimSpace(r.Category),
		IsActive:    true,
	}
}

// ═══════════════════════════════════════════════
// UpdateProductRequest
// ═══════════════════════════════════════════════

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	Category    *string  `json:"category,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

func (r *UpdateProductRequest) Validate() error {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if r.Price != nil && *r.Price < 0 {
		return errors.New("price must be >= 0")
	}
	if r.Stock != nil && *r.Stock < 0 {
		return errors.New("stock must be >= 0")
	}
	return nil
}

func (r *UpdateProductRequest) ToMap() map[string]any {
	m := make(map[string]any)
	if r.Name != nil {
		m["name"] = strings.TrimSpace(*r.Name)
	}
	if r.Description != nil {
		m["description"] = strings.TrimSpace(*r.Description)
	}
	if r.Price != nil {
		m["price"] = *r.Price
	}
	if r.Stock != nil {
		m["stock"] = *r.Stock
	}
	if r.Category != nil {
		m["category"] = strings.TrimSpace(*r.Category)
	}
	if r.IsActive != nil {
		m["is_active"] = *r.IsActive
	}
	return m
}

// ═══════════════════════════════════════════════
// ProductResponse
// ═══════════════════════════════════════════════

type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func FromProduct(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}
