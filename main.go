package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Staring proxy")

	// Init proxy
	proxy := Proxy()
	proxy.Verbose = true

	// Listen on proxy port
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
