package main

import (
	"log"
	"net/http"
)

func main() {
	storage := NewMemStorage()

	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", HandleStorage(storage))
	http.Handle("/stats/", HandleStats(storage))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
