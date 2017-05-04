package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

var (
	server = flag.Bool("s", false, "If set runs in server mode. Otherwise runs in runClient mode.")
	port   = flag.String("p", "9123", "Sets the port to use")
)

func main() {
	flag.Parse()
	if *server {
		runServer()
	} else {
		runClient()
	}
}

func runClient() {
	for i := 0; i < 100; i++ {
		go send()
	}
	time.Sleep(time.Hour)
}

func send() {
	for {
		resp, err := http.Get("http://localhost:" + *port)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}

func runServer() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
