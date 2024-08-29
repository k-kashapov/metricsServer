package main

import (
	"fmt"
	"log"
	"time"
	"runtime"
	"net/http"
	"math/rand"
	"io"
	"os"
	"reflect"
)

var stats = [...]string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", 
						"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", 
						"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", 
						"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", 
						"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

func main() {
	var stat runtime.MemStats
	pollInterval := 2 * time.Second
	reportInterval := 5 * pollInterval
	var timePassed time.Duration = reportInterval

	client := &http.Client{}

	for true {
		runtime.ReadMemStats(&stat)
		time.Sleep(pollInterval)

		timePassed += pollInterval
		if timePassed >= reportInterval {
			timePassed -= reportInterval

			url := fmt.Sprint("http://localhost:8080/update/counter/", "pollCount/", "1")
			// fmt.Println("sending url =", url)
			response, err := client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			// io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			randomValue := rand.Int()
			url = fmt.Sprint("http://localhost:8080/update/gauge/", "randomValue/", randomValue)

			// fmt.Println("sending url =", url)
			response, err = client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			for _, name := range stats {
				val := reflect.ValueOf(stat).FieldByName(name)
				url = fmt.Sprint("http://localhost:8080/update/gauge/", name, "/", val)

				// fmt.Println("sending url =", url)
				response, err = client.Post(url, "text/plain", nil)
				if err != nil {
					log.Fatal("Unable to Post: ", err)
				}

				response.Body.Close()
			}
		}
	}
}
