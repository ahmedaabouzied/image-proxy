package main

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

// Proxy returns a proxy instance
func Proxy() *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()

	// Go proxies always break websites because they canonicalize HTTP headers by default.
	// While this is a good standard, some other languages and server-client systems don't follow
	// this standard.
	// Therefore, Go proxies break websites.
	// In any go proxy, this auto canonicalization of headers should be disabled.
	proxy.PreventCanonicalization = true
	proxy.KeepDestinationHeaders = true
	proxy.KeepHeader = true

	// Intercept response
	proxy.OnResponse().DoFunc(
		// Add X-Image-Proxy header
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			resp.Header.Add("X-Image-Proxy", "1")
			return resp
		},
	)

	return proxy
}
