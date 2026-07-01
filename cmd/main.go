package main

import (
	"fmt"
	"net/http"

	"github.com/Xenios7/Trade-executor/internal/api"
	"github.com/Xenios7/Trade-executor/internal/kafka"
	"github.com/Xenios7/Trade-executor/internal/service"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {

	p, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	producer := kafka.NewKafkaProducer(p, "trade-orders")
	svc := service.NewOrderService(producer, nil)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	fmt.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}