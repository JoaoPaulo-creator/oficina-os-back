// internal/repository/memory/order.go
package memory

import (
	"context"
	"sync"

	"oficina-os/internal/domain"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]domain.Order)}
}

func (r *OrderRepository) Create(_ context.Context, o domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID] = o
	return nil
}

func (r *OrderRepository) List(_ context.Context) ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Order, 0, len(r.orders))
	for _, o := range r.orders {
		out = append(out, o)
	}
	return out, nil
}

func (r *OrderRepository) GetByID(_ context.Context, id string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	return o, nil
}

func (r *OrderRepository) Update(_ context.Context, o domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return domain.ErrOrderNotFound
	}
	r.orders[o.ID] = o
	return nil
}

func (r *OrderRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[id]; !ok {
		return domain.ErrOrderNotFound
	}
	delete(r.orders, id)
	return nil
}
