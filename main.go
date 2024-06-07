package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"redis-epoll/config"
	"redis-epoll/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the server")
	flag.IntVar(&config.Port, "port", 7379, "port for the server")
	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("server started \u2684\uFE0E")

	var sigs chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(1)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, sigs)

	go func() {
        http.ListenAndServe("localhost:6060", nil) // Start pprof server
    }()

	wg.Wait()
}
