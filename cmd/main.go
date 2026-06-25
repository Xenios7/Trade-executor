package main

import (
	"fmt"
	"net/http"

	"github.com/Xenios7/Trade-executor/internal/api"
)

func main() {
	h := api.NewHandler()
	r := api.NewRouter(h)

	fmt.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("HTTP server error:", err)
	}
}