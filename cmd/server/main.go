package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beto-codes/url-shortener/internal/handler"
	"github.com/beto-codes/url-shortener/internal/shortener"
	"github.com/beto-codes/url-shortener/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	service := shortener.NewService(store)

	baseURL := "http://localhost:8080"
	h := handler.NewHandler(service, baseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.Shorten)
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/", h.Redirect)

	port := ":8080"
	fmt.Printf("Server starting on %s\n", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
