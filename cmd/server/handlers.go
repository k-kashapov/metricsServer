package main

import (
	"log"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
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

func HandleStorage(storage MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Request type is not POST")
			return
		}

		trimmed := strings.TrimPrefix(r.URL.Path, "/update/")
		updates := strings.Split(trimmed, "/")

		if len(updates) < 3 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Insufficient arguments")
			return
		}

		name := updates[1]

		switch updates[0] {
		case "gauge":
			val, err := strconv.ParseFloat(updates[2], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Invalid argument: %s\n", updates[2])
				return
			}

			storage.UpdateGauge(name, val)
			fmt.Fprintln(w, name, "is set to", val)

		case "counter":
			val, err := strconv.ParseInt(updates[2], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Invalid argument: %s\n", updates[2])
				return
			}

			storage.UpdateCounter(name, val)
			fmt.Fprintln(w, name, "is set to", storage.Counters[name])

		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid tpye: %s\n", updates[0])
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
