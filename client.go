package main

/*
#include <stdlib.h>
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"sync"
	"time"

	"github.com/codahale/hdrhistogram"
)

func runClientTrace(numSenders, requests int) {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	runClient(numSenders, requests)
}

func runClient(numSenders, requests int) {
	fmt.Printf("Client about to send %d requests * %d senders\n", requests, numSenders)
	histChan := make(chan (*hdrhistogram.Histogram), numSenders)
	wg := &sync.WaitGroup{}
	for i := 0; i < numSenders; i++ {
		wg.Add(1)
		go send(requests, wg, histChan)
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
	for _, f := range []float64{0.01, 1, 25, 50, 75, 90, 95, 99, 100} {
		v := hist.ValueAtQuantile(f)
		fmt.Printf("%3.0f: %8dms\n", f, v/int64(time.Millisecond))
	}
}

func send(requests int, wg *sync.WaitGroup, ch chan *hdrhistogram.Histogram) {
	hist := newHistogram()
	for i := 0; i < requests; i++ {
		before := time.Now()
		resp, err := http.Get("http://localhost:" + *port + "/test")
		if err != nil {
			log.Fatal(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		hist.RecordValue(int64(time.Since(before)))
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	ch <- hist
	wg.Done()
}

func newHistogram() *hdrhistogram.Histogram {
	return hdrhistogram.New(0, int64(time.Minute), 5)
}
