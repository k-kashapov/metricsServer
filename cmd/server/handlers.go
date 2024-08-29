package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func HandleStats(storage MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(storage)
		if err != nil {
			log.Print("Could not marshall storage")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func HandleGet(storage MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vtype := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch vtype {
		case "gauge":
			val, ok := storage.Gauges[name]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Not found:", name)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, val)
		case "counter":
			val, ok := storage.Counters[name]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Not found:", name)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, val)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Invalid type:", vtype)
			return
		}
	}
}

func HandleStorage(storage MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vtype := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		switch vtype {
		case "gauge":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid argument:", value)
				return
			}

			storage.UpdateGauge(name, val)
			fmt.Fprintln(w, name, "is set to", val)

		case "counter":
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid argument:", value)
				return
			}

			storage.UpdateCounter(name, val)
			fmt.Fprintln(w, name, "is set to", storage.Counters[name])

		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid type:", vtype)
			return
		}

		// fmt.Print(w, "gauges:\n")
		// for key, val := range storage.Gauges {
		// 	fmt.Printf("\t%s = %g\n", key, val)
		// }

		// fmt.Print(w, "counters:\n")
		// for key, val := range storage.Counters {
		// 	fmt.Printf("\t%s = %d\n", key, val)
		// }
	}
}
