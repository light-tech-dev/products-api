package item

import (
	"context"

	"github.com/abdallah-elngar/gormx"
)

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Service{ctx: ctx}
}

func (s *Service) Query() *gormx.QuerySet[Item] {
	return gormx.New[Item]().WithContext(s.ctx)
}

func (s *Service) Create(productID uint, req *CreateItemRequest) (*ItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	item := &Item{
		ProductID: productID,
		Name:      req.Name,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Notes:     req.Notes,
	}

	if err := s.Query().Create(item); err != nil {
		return nil, err
	}

	return FromItem(item), nil
}

func (s *Service) ListByProduct(productID uint) ([]ItemResponse, error) {
	items, err := s.Query().
		Filter("product_id", productID).
		OrderBy("id").
		All()
	if err != nil {
		return nil, err
	}

	out := make([]ItemResponse, len(items))
	for i, item := range items {
		out[i] = *FromItem(&item)
	}
	return out, nil
}

func (s *Service) Delete(id uint) error {
	return s.Query().Delete(id)
}
