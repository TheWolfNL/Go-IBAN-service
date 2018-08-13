package main

import (
	"log"
	"os"
)

func main() {
	loadConfig()
	startServer()
}

func loadConfig() {
	if os.Getenv("SERVICE_HOST") != "" {
		serverConfig.host = os.Getenv("SERVICE_HOST")
	}

	if os.Getenv("SERVICE_PORT") != "" {
		serverConfig.port = os.Getenv("SERVICE_PORT")
	}

	if os.Getenv("OPTIONS_SANITZE") != "" {
		serverConfig.sanitize = true
	}

	log.Printf("Config: %+v\n", serverConfig)
}
