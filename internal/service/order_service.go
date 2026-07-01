package service

import (
	"fmt"
	"time"

	"github.com/Xenios7/Trade-executor/internal/domain"
)

type Producer interface {
	Publish(order domain.Order) error
}

type Repository interface {
	Save(domain.Order) error
	Cache(domain.Order) error
}

type OrderService struct {
	producer Producer
	repo Repository	
}

func NewOrderService(p Producer, r Repository) *OrderService {
	return &OrderService{
		producer: p,
		repo: r,
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
	if s.repo != nil {
		if err := s.repo.Save(order); err != nil {
			return err
		}
		if err := s.repo.Cache(order); err != nil {
			return err
		}
	}

    return nil
}

// Stubs for now 
func (s *OrderService) GetOrder(id string) (domain.Order, error) {
	return domain.Order{}, nil
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
	return []domain.Order{}, nil
}	
