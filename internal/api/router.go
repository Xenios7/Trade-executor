package api

import (
	"net/http"
	
)

func NewRouter(h *Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.PublishHandler)
	mux.HandleFunc("GET /orders", h.GetAllOrdersHandler)
	mux.HandleFunc("GET /orders/{id}", h.GetOrderHandler)

	return mux	
}