package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	proxy := Proxy()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var not found")
	}
	fmt.Printf("Proxy server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), proxy))
}
