FROM golang:1.26.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev librdkafka-dev pkgconfig
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -o trade-executor ./cmd

FROM alpine:3.19
RUN apk add --no-cache librdkafka
WORKDIR /app
COPY --from=builder /app/trade-executor .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./trade-executor"]