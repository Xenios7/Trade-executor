package api

import (
	"net/http"
)

func NewRouter(h *Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /orders", h.PublishHandler)
	mux.HandleFunc("GET /orders", h.ConsumeHandler)

	
}