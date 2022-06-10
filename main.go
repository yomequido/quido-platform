// main.go

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/yomequido/quido-platform/platform/authenticator"
	"github.com/yomequido/quido-platform/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	log.Print("Server listening on http://localhost:" + port)
	if err := http.ListenAndServe("0.0.0.0:"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
