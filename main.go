package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const portEnvKey = "PORT"

func main() {
	proxy := Proxy()

	port := os.Getenv(portEnvKey)
	if port == "" {
		log.Fatalf("%s env var must be set", portEnvKey)
	}
	fmt.Printf("Proxy server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), proxy))
}
