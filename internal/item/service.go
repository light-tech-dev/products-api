package item

import (
	"context"

	"github.com/abdallah-elngar/gormx"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Query(ctx context.Context) *gormx.QuerySet[Item] {
	return gormx.New[Item]().WithContext(ctx)
}

func (s *Service) Create(ctx context.Context, productID uint, req *CreateItemRequest) (*ItemResponse, error) {
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

	if err := s.Query(ctx).Create(item); err != nil {
		return nil, err
	}

	return FromItem(item), nil
}

func (s *Service) ListByProduct(ctx context.Context, productID uint) ([]ItemResponse, error) {
	items, err := s.Query(ctx).
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

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.Query(ctx).Delete(id)
}
