package main

import (
	"fmt"
	"flag"
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
	addr := flag.String("a", "localhost:8080", "endpoint address")
	reportSec := flag.Int("r", 10, "interval in seconds between reports to the server")
	pollSec := flag.Int("p", 2, "interval in seconds between polling the stats")

	flag.Parse()

	reportInterval := time.Duration(*reportSec) * time.Second
	pollInterval := time.Duration(*pollSec) * time.Second

	var stat runtime.MemStats
	var timePassed time.Duration = reportInterval

	client := &http.Client{}

	for true {
		runtime.ReadMemStats(&stat)
		time.Sleep(pollInterval)

		fmt.Println("Sleep over")

		timePassed += pollInterval
		if timePassed >= reportInterval {
			timePassed -= reportInterval

			fmt.Println("Take over")
			url := fmt.Sprint("http://", *addr, "/update/counter/", "pollCount/", "1")
			// fmt.Println("sending url =", url)
			response, err := client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			// io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			randomValue := rand.Int()
			url = fmt.Sprint("http://", *addr, "/update/gauge/", "randomValue/", randomValue)

			// fmt.Println("sending url =", url)
			response, err = client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			for _, name := range stats {
				val := reflect.ValueOf(stat).FieldByName(name)
				url = fmt.Sprint("http://", *addr, "/update/gauge/", name, "/", val)

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
