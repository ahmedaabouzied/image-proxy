package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const portEnvKey = "PORT"

func main() {
	go func() {
		interceptor := NewInterceptor()

		port := os.Getenv(portEnvKey)
		if port == "" {
			log.Fatalf("%s env var must be set", portEnvKey)
		}

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), interceptor))
	}()

	startUI()
}
