package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"
)

var stats = [...]string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
	"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys",
	"PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

type Config struct {
	Addr    string        `env:"ADDRESS"`
	Rep  time.Duration `env:"REPORT_INTERVAL"`
	Poll time.Duration `env:"POLL_INTERVAL"`
}

func main() {
	addrPtr := flag.String("a", "localhost:8080", "endpoint address")
	reportSec := flag.Int("r", 10, "interval in seconds between reports to the server")
	pollSec := flag.Int("p", 2, "interval in seconds between polling the stats")

	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if cfg.Addr == "" {
		cfg.Addr = *addrPtr
	}

	if cfg.Rep == 0 {
		cfg.Rep = time.Duration(*reportSec)
	}

	if cfg.Poll == 0 {
		cfg.Poll = time.Duration(*pollSec)
	}

	reportInterval := cfg.Rep * time.Second
	pollInterval := cfg.Poll * time.Second

	log.Printf("Running agent with config: addr=%s, reportInterval=%s, pollInterval=%s", cfg.Addr, reportInterval, pollInterval)

	var stat runtime.MemStats
	var timePassed time.Duration

	client := &http.Client{}

	for {
		runtime.ReadMemStats(&stat)
		time.Sleep(pollInterval)

		timePassed += pollInterval
		if timePassed >= reportInterval {
			timePassed -= reportInterval

			url := fmt.Sprint("http://", cfg.Addr, "/update/counter/", "pollCount/", "1")
			// fmt.Println("sending url =", url)
			response, err := client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			// io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			randomValue := rand.Int()
			url = fmt.Sprint("http://", cfg.Addr, "/update/gauge/", "randomValue/", randomValue)

			// fmt.Println("sending url =", url)
			response, err = client.Post(url, "text/plain", nil)
			if err != nil {
				log.Fatal("Unable to Post: ", err)
			}

			io.Copy(os.Stdout, response.Body)
			response.Body.Close()

			for _, name := range stats {
				val := reflect.ValueOf(stat).FieldByName(name)
				url = fmt.Sprint("http://", cfg.Addr, "/update/gauge/", name, "/", val)

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
