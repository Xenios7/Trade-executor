package main

import (
	"fmt"
	"net/http"

	"github.com/Xenios7/Trade-executor/internal/api"
	"github.com/Xenios7/Trade-executor/internal/service"
)

func main() {
	svc := service.NewOrderService(nil, nil)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	fmt.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}