package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":20008"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./_tmp"
	}

	os.MkdirAll(dataDir, 0755)

	registerHandlers(dataDir)

	log.Printf("Server listening on %s", listenAddr)
	log.Printf("Data directory: %s", dataDir)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
