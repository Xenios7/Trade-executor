package domain

import "time"

type Order struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	
	CreatedAt   time.Time  `json:"created_at"`
	ExecutedAt  *time.Time `json:"executed_at"`
	Asset       string     `json:"asset"`
	Side        string     `json:"side"`
	Quantity    float64    `json:"quantity"`
	Price       float64    `json:"price"`
}