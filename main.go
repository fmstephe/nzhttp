package main

import (
	"flag"
	"time"
)

var (
	server      = flag.Bool("server", false, "If set runs in server mode.")
	serverPause = flag.Duration("pause", 0, "The time for the server to pause before responding")
	poolServer  = flag.Bool("serverPool", false, "If set runs in worker pool server mode.")
	poolSize    = flag.Int("poolSize", 10, "The size of the worker pool for the pool server")
	port        = flag.String("port", "9123", "Sets the port to use")
	useC        = flag.Bool("useC", false, "If set the server will sleep using cGo")
)
var (
	client   = flag.Bool("client", false, "If set runs in client mode.")
	requests = flag.Int("requests", 1000, "Number of requests sent by each sender in client mode")
	senders  = flag.Int("senders", 100, "Number of senders to use in client mode")
)

func main() {
	flag.Parse()
	if *server {
		runServer(*port, *serverPause, *useC)
	} else if *poolServer {
		runServerPool(*port, *serverPause, *useC, *poolSize)
	} else if *client {
		for i := 0; i < 10; i++ {
			runClient(*senders, *requests)
			time.Sleep(time.Second)
		}
	}
}
