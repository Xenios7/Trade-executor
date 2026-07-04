package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Xenios7/Trade-executor/internal/api"
	"github.com/Xenios7/Trade-executor/internal/kafka"
	"github.com/Xenios7/Trade-executor/internal/repository"
	"github.com/Xenios7/Trade-executor/internal/service"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	goredis "github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"

)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	
	// 0. Environment variables
	// Load .env file for local development
    godotenv.Load()
	postgresDSN := getEnv("POSTGRES_DSN", "")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")

	if postgresDSN == "" {
		panic("POSTGRES_DSN environment variable not set")
	}

	// 1. Postgres, connect and run migrations
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		panic(err)
	}
	m, err := migrate.New("file://migrations", postgresDSN)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	// 2. Redis, connect
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: redisAddr,
	})

	// 3. Repositories
	store := repository.NewPostgresRepository(db)
	cache := repository.NewRedisRepository(redisClient)

	// 4. Kafka producer
	p, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
	})
	if err != nil {
		panic(err)
	}
	producer := kafka.NewKafkaProducer(p, "trade-orders")

	// 5. Service and HTTP server
	svc := service.NewOrderService(producer, store, cache)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	// 6. Kafka consumer, runs in background goroutine
	c, err := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          "trade-executor",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	consumer := kafka.NewKafkaConsumer(c, svc)
	go consumer.Start()

	// 7. Start HTTP server, blocks main goroutine
	fmt.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}