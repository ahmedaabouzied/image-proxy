package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	proxy := Proxy()

	port := 8080
	fmt.Printf("Proxy server running on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), proxy))
}
