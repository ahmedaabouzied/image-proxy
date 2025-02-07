package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	fmt.Println("Staring proxy")
	proxy := Proxy()
	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func Proxy() *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnResponse().DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			resp.Header.Add("X-Image-Proxy", "1")
			return resp
		},
	)

	return proxy
}
