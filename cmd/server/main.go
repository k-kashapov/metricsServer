package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func NewRouter(storage MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", logHandler(HandleStorage(storage)))
	r.Get("/value/{type}/{name}", logHandler(HandleGet(storage)))
	r.Get("/", logHandler(HandleStats(storage)))

	return r
}

func main() {
	addrPtr := flag.String("a", "localhost:8080", "endpoint address")

	var addr string
	addr, ok := os.LookupEnv("ADDRESS")
	if !ok {
		flag.Parse()
		addr = *addrPtr
	}

	log.Println("Running server at", addr)

	storage := NewMemStorage()
	r := NewRouter(storage)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
