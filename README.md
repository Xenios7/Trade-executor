# Trade-executor

An event-driven cryptocurrency trade execution pipeline built in Go. Trade orders are submitted via HTTP API, published to Apache Kafka, processed asynchronously by a consumer group, persisted in PostgreSQL, and cached in Redis. Deployed on AWS ECS.

## Architecture

- **API Server** — accepts and validates incoming trade orders (buy/sell, asset, quantity, price)
- **Kafka Producer** — publishes validated orders as events to the `trade-orders` topic
- **Kafka Consumer** — consumes events, simulates execution, updates order status (PENDING → FILLED / REJECTED)
- **Redis** — caches full order state for low-latency lookups
- **PostgreSQL** — persists all orders and their final execution state
- **AWS ECS** — container orchestration for production deployment

```
POST /orders → validate → produce to Kafka → consumer processes →
FILLED / REJECTED → persist to Postgres → cache in Redis → GET /orders/{id} reads from Redis
```

## Tech Stack

Go, Apache Kafka, Redis, PostgreSQL, Docker, AWS ECS

## Milestones

| | Milestone | Status |
|---|---|---|
| M1 | HTTP API server — accept, validate, and store trade orders | ✅ Complete |
| M2 | Kafka producer — publish validated orders to Kafka topic | ✅ Complete |
| M3 | Kafka consumer — process orders asynchronously, update status | ✅ Complete |
| M4 | Redis caching — cache full order state for fast lookups | ✅ Complete |
| M5 | PostgreSQL persistence — persist all orders and execution results | ✅ Complete |
| M6 | Docker Compose — run full stack locally | ✅ Complete |
| M7 | AWS ECS deployment — deploy containerized service to AWS | ⬜ Not Started |

## Order Flow

A trade order moves through the following states:

```
PENDING → FILLED   
         → REJECTED  
```

A new order is stored as `PENDING`. The Kafka consumer processes it and updates the status to `FILLED` or `REJECTED` based on order size. The result is persisted to PostgreSQL and cached in Redis.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | /orders | Submit a new trade order |
| GET | /orders/{id} | Retrieve order by ID |
| GET | /orders | List all orders |

## Example

**Submit Order**
```
POST /orders
Content-Type: application/json

{
  "asset": "BTC/USD",
  "side": "buy",
  "quantity": 0.5,
  "price": 65000.00
}
```

**Response (202 Accepted)**
```json
{
  "id": "70df8905-54c1-43ca-9cf8-b134587ae732",
  "asset": "BTC/USD",
  "side": "buy",
  "quantity": 0.5,
  "price": 65000,
  "status": "PENDING",
  "created_at": "2026-07-04T08:59:14.128126Z",
  "executed_at": null
}
```

**Check Order**
```
GET /orders/70df8905-54c1-43ca-9cf8-b134587ae732
```
```json
{
  "id": "70df8905-54c1-43ca-9cf8-b134587ae732",
  "asset": "BTC/USD",
  "side": "buy",
  "quantity": 0.5,
  "price": 65000,
  "status": "FILLED",
  "created_at": "2026-07-04T08:59:14.128126Z",
  "executed_at": "2026-07-04T08:59:17.888758Z"
}
```

## Running Locally

```bash
docker compose up -d
go run cmd/main.go
```

## Author

Xenios Gerolemou — [LinkedIn](https://www.linkedin.com/in/xenios-gerolemou-594086202/) · [Portfolio](https://xenios7.github.io/portfolio/)