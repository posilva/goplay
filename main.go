package main

import (
	"github.com/posilva/goplay/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	host     = "localhost"
	port     = "3333"
	connType = "tcp"
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-sigs
		done <- true
	}()

	srv := server.New()
	err := srv.Listen()
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	os.Exit(srv.Start())
}
