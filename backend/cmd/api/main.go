package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/adapters/server"
)

func main() {

	apiServer := server.NewServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go server.GracefulShutdown(apiServer, done)

	err := apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
