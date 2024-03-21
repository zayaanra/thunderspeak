package main

import (
	"os"
	"os/signal"
	"syscall"

	https "github.com/zayaanra/thunderspeak/https/server"
)

func main() {
	// Start the HTTP server on port 8080
	s := https.NewServer()

	if err := s.Open(); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	if err := s.Close(); err != nil {
		panic(err)
	}

}
