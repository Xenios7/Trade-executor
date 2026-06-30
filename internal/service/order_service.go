package service

import (
	"github.com/Xenios7/Trade-executor/internal/domain"
)

type Producer interface {
	Publish(order domain.Order) error
}

type Repository interface {
	//
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

func (s *OrderService) PlaceOrder(order domain.Order) error {
	if s.producer != nil {
		if err := s.producer.Publish(order); err != nil {
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
