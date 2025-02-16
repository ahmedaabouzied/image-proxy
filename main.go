package main

import (
    "flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
    port := flag.String("port", "8080", "Port to run the server on")
    flag.Parse()

	interceptor := NewInterceptor()

    log.Printf("Starting server on port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), interceptor))
}
