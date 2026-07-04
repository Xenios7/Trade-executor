FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o trade-executor ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/trade-executor .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./trade-executor"]