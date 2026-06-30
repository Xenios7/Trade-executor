package kafka

import (
	"encoding/json"

	"github.com/Xenios7/Trade-executor/internal/domain"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
    producer *ckafka.Producer
    topic    string
}

func NewKafkaProducer(prod *ckafka.Producer, t string) *KafkaProducer {
	return &KafkaProducer{
		producer: prod,
		topic: t,
	}
}

// Publish serializes an Order and sends it to the Kafka topic.
// The Asset field is used as the message key so that all orders
// for the same asset (e.g. BTC/USD) always land on the same partition,
// guaranteeing ordering per asset.
func (p *KafkaProducer) Publish(order domain.Order) error {

	// Serialize the Order struct to JSON bytes
	// e.g. {"id":"abc","asset":"BTC/USD","side":"buy",...}
	value, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return p.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &p.topic,       // which topic to publish to (trade-orders)
			Partition: ckafka.PartitionAny, // let Kafka pick partition based on key
		},
		Key:   []byte(order.Asset), // partition key — BTC/USD always → same partition
		Value: value,               // the serialized order as JSON bytes
	}, nil) // nil = no delivery report channel (fire and forget)
}