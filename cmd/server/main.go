package main

import (
	"flag"
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
	addr := flag.String("a", "localhost:8080", "endpoint address")
	flag.Parse()

	storage := NewMemStorage()
	r := NewRouter(storage)

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
