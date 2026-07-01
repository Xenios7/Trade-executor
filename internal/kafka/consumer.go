package kafka

import (
	"encoding/json"

	"github.com/Xenios7/Trade-executor/internal/domain"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type OrderService interface {
	ProcessOrder(order domain.Order) error
 }

type KafkaConsumer struct {
	consumer *ckafka.Consumer
	service OrderService 
}

func NewKafkaConsumer(c *ckafka.Consumer, s OrderService) *KafkaConsumer {
	return &KafkaConsumer{
		consumer: c,
		service: s,
	}
}

func (c *KafkaConsumer) Start() {
    // Subscribe to the trade-orders topic
    // slice because you can subscribe to multiple topics at once
    c.consumer.SubscribeTopics([]string{"trade-orders"}, nil)

    // Loop forever — consumer never stops, keeps waiting for new messages
    for {
        // Block until a message arrives (-1 = no timeout, wait indefinitely)
        msg, err := c.consumer.ReadMessage(-1)
        if err != nil {
            // Log and continue — don't crash on a single bad message
            continue
        }

        // Deserialize JSON bytes back into a domain.Order struct
        var order domain.Order
        if err := json.Unmarshal(msg.Value, &order); err != nil {
            continue
        }

        // Hand off to service layer for processing (FILLED/REJECTED logic)
        c.service.ProcessOrder(order)
    }
}