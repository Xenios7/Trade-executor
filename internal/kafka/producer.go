package kafka

type Producer interface {
	Publish(order Order) error
}

