package service

import (
	"fmt"
	"time"

	"github.com/Xenios7/Trade-executor/internal/domain"
)

type Producer interface {
	Publish(order domain.Order) error
}

type Store interface {
    Save(order domain.Order) error
    GetByID(id string) (domain.Order, error)
	GetAll() ([]domain.Order, error)
}	

type Cache interface {
    Cache(order domain.Order) error
    GetCached(id string) (domain.Order, error)
}

type OrderService struct {
	producer Producer
	store Store
	cache Cache
}

func NewOrderService(p Producer, s Store, c Cache) *OrderService {
	return &OrderService{
		producer: p,
		store: s,
		cache: c,
	}
}
//PlaceOrder
func (s *OrderService) PlaceOrder(order domain.Order) error {
	if s.producer != nil {
		if err := s.producer.Publish(order); err != nil {
			return err
		}
	}
	return nil
}

//ProcessOrder
func (s *OrderService) ProcessOrder(order domain.Order) error {
    // decide FILLED or REJECTED based on order size
    if order.Price * order.Quantity > 1000000 {
        order.Status = "REJECTED"
    } else {
        order.Status = "FILLED"
    }
    fmt.Printf("Order %s processed: %s\n", order.ID, order.Status)
	
    // set execution timestamp
    now := time.Now().UTC()
    order.ExecutedAt = &now

    // persist and cache
	if s.store != nil {
		if err := s.store.Save(order); err != nil {
			return err
		}
	}
	if s.cache != nil {
		if err := s.cache.Cache(order); err != nil {
			return err
		}
	}
    return nil
}

func (s *OrderService) GetOrder(id string) (domain.Order, error) {
    // Check Redis first
    if s.cache != nil {
        order, err := s.cache.GetCached(id)
        if err == nil {
            return order, nil
        }
    }
    // Fall back to Postgres
    if s.store != nil {
        return s.store.GetByID(id)
    }
    return domain.Order{}, nil
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
    if s.store != nil {
        return s.store.GetAll()
    }
    return []domain.Order{}, nil
}