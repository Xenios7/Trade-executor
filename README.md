# Trade-executor

An event-driven cryptocurrency trade execution pipeline built in Go. Trade orders are submitted via HTTP API, published to Apache Kafka, processed asynchronously by a consumer group, and persisted in PostgreSQL with results cached in Redis. Deployed on AWS ECS.

## Architecture

- **API Server** — accepts and validates incoming trade orders (buy/sell, asset, quantity, price)
- **Kafka Producer** — publishes validated orders as events to the `trade-orders` topic
- **Kafka Consumer** — consumes events, simulates execution, updates order status (PENDING → FILLED / REJECTED)
- **Redis** — caches order state for low-latency status lookups
- **PostgreSQL** — persists all orders and their final execution state
- **AWS ECS** — container orchestration for production deployment

```
POST /orders → validate → produce to Kafka → consumer processes →
update status → write to Redis + Postgres → GET /orders/{id} reads from Redis
```

## Tech Stack

Go, Apache Kafka, Redis, PostgreSQL, Docker, AWS ECS

## Milestones

| | Milestone | Status |
|---|---|---|
| M1 | HTTP API server — accept, validate, and store trade orders | ✅ Complete |
| M2 | Kafka producer — publish validated orders to Kafka topic | 🔄 In Progress |
| M3 | Kafka consumer — process orders asynchronously, update status | ⬜ Not Started |
| M4 | Redis caching — cache order state for fast lookups | ⬜ Not Started |
| M5 | PostgreSQL persistence — persist all orders and execution results | ⬜ Not Started |
| M6 | Docker Compose — run full stack locally | ⬜ Not Started |
| M7 | AWS ECS deployment — deploy containerized service to AWS | ⬜ Not Started |
## Order Flow

A trade order moves through the following states:

```
PENDING → FILLED
         → REJECTED
```

A new order is stored as `PENDING`. The Kafka consumer processes it and updates the status to `FILLED` or `REJECTED` based on execution logic.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | /orders | Submit a new trade order |
| GET | /orders/{id} | Retrieve order status by ID |
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
  "id": "a1b2c3",
  "asset": "BTC/USD",
  "side": "buy",
  "quantity": 0.5,
  "price": 65000.00,
  "status": "PENDING",
  "created_at": "2026-06-23T20:00:00Z"
}
```

**Check Status**
```
GET /orders/a1b2c3
```
```json
{
  "id": "a1b2c3",
  "status": "FILLED",
  "executed_at": "2026-06-23T20:00:01Z"
}
```

## Running Locally

```bash
docker compose up --build
```

## Author

Xenios Gerolemou — [LinkedIn](https://www.linkedin.com/in/xenios-gerolemou-594086202/) · [Portfolio](https://xenios7.github.io/portfolio/)
