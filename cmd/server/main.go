package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func NewRouter(storage MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", HandleStorage(storage))
	r.Get("/value/{type}/{name}", HandleGet(storage))
	r.Get("/", HandleStats(storage))

	return r
}

func main() {
	storage := NewMemStorage()
	r := NewRouter(storage)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
