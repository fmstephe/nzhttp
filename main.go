package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codahale/hdrhistogram"
)

var (
	server   = flag.Bool("s", false, "If set runs in server mode. Otherwise runs in runClient mode.")
	port     = flag.String("p", "9123", "Sets the port to use")
	requests = flag.Int("n", 1000, "Number of requests sent by each sender in client mode")
	senders  = flag.Int("m", 100, "Number of senders to use in client mode")
)

func main() {
	flag.Parse()
	if *server {
		runServer()
	} else {
		runClient()
	}
}

func runServer() {
	fmt.Printf("Server starting on port %s\n", *port)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func runClient() {
	numSenders := *senders

	fmt.Printf("Client about to send %d requests * %d senders\n", *requests, numSenders)
	histChan := make(chan (*hdrhistogram.Histogram), numSenders)
	wg := &sync.WaitGroup{}
	for i := 0; i < numSenders; i++ {
		wg.Add(1)
		go send(wg, histChan)
	}
	wg.Wait()
	close(histChan)

	// combine results from workers into one histogram
	hist := newHistogram()
	var dropped int64
	for curHist := range histChan {
		dropped += hist.Merge(curHist)
	}

	// print some results
	fmt.Printf("Merging dropped %d samples\n", dropped)
	for _, f := range []float64{1, 25, 50, 75, 90, 95, 99, 100} {
		v := hist.ValueAtQuantile(f)
		fmt.Printf("%3.0f: %8dms\n", f, v/int64(time.Millisecond))
	}

}

func send(wg *sync.WaitGroup, ch chan *hdrhistogram.Histogram) {
	hist := newHistogram()
	for i := 0; i < *requests; i++ {
		before := time.Now()
		resp, err := http.Get("http://localhost:" + *port)
		hist.RecordValue(int64(time.Since(before)))

		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
	ch <- hist
	wg.Done()
}

func newHistogram() *hdrhistogram.Histogram {
	return hdrhistogram.New(0, int64(time.Minute), 5)
}
