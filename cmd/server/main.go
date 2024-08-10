package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64

type MemStorage struct {
	gauges   map[string]gauge
	counters map[string]counter
}

var storage MemStorage

func (st *MemStorage) updateGauge(name string, val float64) {
	st.gauges[name] = gauge(val)
}

func (st *MemStorage) updateCounter(name string, val int64) {
	st.counters[name] += counter(val)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		trimmed := strings.TrimPrefix(r.URL.Path, "/update/")
		updates := strings.Split(trimmed, "/")

		if len(updates) < 3 {
			// log.Print("Insufficient arguments\n")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		name := updates[1]

		switch updates[0] {
		case "gauge":
			val, err := strconv.ParseFloat(updates[2], 64)
			if err != nil {
				// log.Print("Invalid argument ", updates[2], "\n")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			storage.updateGauge(name, val)
		case "counter":
			val, err := strconv.ParseInt(updates[2], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				// log.Print("Invalid argument ", updates[2], "\n")
				return
			}

			storage.updateCounter(name, val)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		// fmt.Fprint(w, "gauges:\n")
		// for key, val := range storage.gauges {
		// 	fmt.Fprintf(w, "\t%s = %g\n", key, val)
		// }

		// fmt.Fprint(w, "counters:\n")
		// for key, val := range storage.counters {
		// 	fmt.Fprintf(w, "\t%s = %d\n", key, val)
		// }
	}
}

func main() {
	storage.gauges = make(map[string]gauge, 0)
	storage.counters = make(map[string]counter, 0)

	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", handleUpdate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
