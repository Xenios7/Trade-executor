package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type PublishRequest struct {
	Asset    string  `json:"asset"`
	Side     string  `json:"side"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderResponse struct {
	ID        string    `json:"id"`
	Asset     string    `json:"asset"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response := OrderResponse{
		ID:        uuid.New().String(),
		Asset:     req.Asset,
		Side:      req.Side,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Status:    "PENDING",
		CreatedAt: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}