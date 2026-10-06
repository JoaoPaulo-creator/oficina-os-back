// internal/service/order.go
package service

import (
	"context"
	"oficina-os/internal/domain"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, o domain.Order) error
	List(ctx context.Context) ([]domain.Order, error)
	GetByID(ctx context.Context, id string) (domain.Order, error)
	Update(ctx context.Context, o domain.Order) error
	Delete(ctx context.Context, id string) error
}

type OrderService struct {
	repo  OrderRepository
	newID func() string
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo, newID: uuid.NewString}
}

func (s *OrderService) Create(ctx context.Context, o domain.Order) (domain.Order, error) {
	if err := o.Validate(); err != nil {
		return domain.Order{}, err
	}
	o.ID = s.newID()
	if err := s.repo.Create(ctx, o); err != nil {
		return domain.Order{}, err
	}
	return o, nil
}

func (s *OrderService) List(ctx context.Context) ([]domain.Order, error) {
	return s.repo.List(ctx)
}

func (s *OrderService) Get(ctx context.Context, id string) (domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) Update(ctx context.Context, id string, o domain.Order) (domain.Order, error) {
	o.ID = id
	if err := o.Validate(); err != nil {
		return domain.Order{}, err
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return domain.Order{}, err
	}
	return o, nil
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
