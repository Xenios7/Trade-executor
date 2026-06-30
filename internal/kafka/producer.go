package kafka

import "github.com/Xenios7/Trade-executor/internal/domain"


type Producer interface {
	Publish(order domain.Order) error
}

