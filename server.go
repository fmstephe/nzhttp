package main

/*
#include <stdlib.h>
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func runServer(port string, pause time.Duration, useC bool) {
	fmt.Printf("Server starting on port %s, pause: %s\n", port, pause.String())
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		s := doServeC(r.URL.Path, pause, useC)
		fmt.Fprint(w, s)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runServerPool(port string, pause time.Duration, useC bool, poolSize int) {
	fmt.Printf("Server starting on port %s, pause: %s, poolSize: %d\n", port, pause.String(), poolSize)
	p := newWorkerPool(poolSize)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		sc := make(chan string, 1)
		f := func() {
			sc <- doServeC(r.URL.Path, pause, useC)
		}
		p.submit(f)
		s := <-sc
		fmt.Fprint(w, s)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func doServeC(path string, pause time.Duration, useC bool) string {
	if useC {
		micros := int(pause / time.Microsecond)
		C.usleep(C.__useconds_t(micros)) // Linux
		//C.usleep(C.useconds_t(micros)) // Darwin
	} else {
		time.Sleep(pause)
	}
	return fmt.Sprintf("Hello, %q", html.EscapeString(path))
}
