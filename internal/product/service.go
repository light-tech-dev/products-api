package product

import (
	"context"
	"errors"
	"strings"

	"github.com/abdallah-elngar/gormx"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrSKUAlreadyExists = errors.New("SKU already exists")
)

// Service handles product business logic.
type Service struct{}

// NewService creates a new ProductService.
func NewService() *Service {
	return &Service{}
}

// Query returns a QuerySet with the given context.
func (s *Service) Query(ctx context.Context) *gormx.QuerySet[Product] {
	return gormx.New[Product]().WithContext(ctx)
}

// Create creates a new product.
func (s *Service) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	exists, err := s.Query(ctx).Filter("sku", strings.ToUpper(req.SKU)).Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSKUAlreadyExists
	}

	product := req.ToModel()
	if err := s.Query(ctx).Create(product); err != nil {
		return nil, err
	}

	return FromProduct(product), nil
}

// GetByID fetches a product by ID.
func (s *Service) GetByID(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.Query(ctx).Get(id)
	if err != nil {
		if gormx.IsNotFound(err) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return FromProduct(product), nil
}

// List returns paginated products with filters.
func (s *Service) List(ctx context.Context, page, perPage int, filters map[string]any) (*gormx.PaginatedResult[ProductResponse], error) {
	q := s.Query(ctx)

	if category, ok := filters["category"].(string); ok && category != "" {
		q = q.Filter("category", category)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		q = q.Filter("is_active", isActive)
	}
	if minPrice, ok := filters["min_price"].(float64); ok {
		q = q.Filter("price__gte", minPrice)
	}
	if maxPrice, ok := filters["max_price"].(float64); ok {
		q = q.Filter("price__lte", maxPrice)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		q = q.Q(gormx.QOr(
			gormx.Contains("name", search),
			gormx.Contains("sku", search),
		))
	}

	result, err := q.OrderBy("-created_at").Paginate(page, perPage)
	if err != nil {
		return nil, err
	}

	items := make([]ProductResponse, len(result.Items))
	for i, p := range result.Items {
		items[i] = *FromProduct(&p)
	}

	return &gormx.PaginatedResult[ProductResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       result.Page,
		PerPage:    result.PerPage,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}, nil
}

// Update updates a product.
func (s *Service) Update(ctx context.Context, id uint, req *UpdateProductRequest) (*ProductResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}

	updates := req.ToMap()
	if len(updates) == 0 {
		return s.GetByID(ctx, id)
	}

	if _, err := s.Query(ctx).Filter("id", id).UpdateMany(updates); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// Delete removes a product.
func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.Query(ctx).Delete(id)
}
