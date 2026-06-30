package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Xenios7/Trade-executor/internal/domain"
	"github.com/google/uuid"
)

type OrderService interface {
    PlaceOrder(order domain.Order) error
    GetOrder(id string) (domain.Order, error)
    GetAllOrders() ([]domain.Order, error)
}

type Handler struct {
	service OrderService
}

func NewHandler(svc OrderService) *Handler {
	return &Handler{
		service: svc,
	}
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
	ExecutedAt  *time.Time `json:"executed_at"`
}


func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	order := domain.Order {
		ID: uuid.New().String(),
		Asset: req.Asset,
		Side: req.Side,
		Quantity: req.Quantity,
		Price: req.Price,
		Status: "PENDING",
		CreatedAt: time.Now().UTC(),
	}
	
	if err := h.service.PlaceOrder(order); err != nil {
		http.Error(w, "failed to place order", http.StatusInternalServerError)
		return
	}

	response := OrderResponse{
		ID:        order.ID,
		Asset:     order.Asset,
		Side:      order.Side,
		Quantity:  order.Quantity,
		Price:     order.Price,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAllOrdersHandler(w http.ResponseWriter, r *http.Request) {

	//kafka producer call
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]OrderResponse{})
}

func (h *Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(OrderResponse{})
}

