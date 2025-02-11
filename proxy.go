package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

const (
	placeholderImageEnvKey = "PLACEHOLDER_IMAGE_PATH"
	publicCertEnvKey       = "PUB_CERT_PATH"
	privateCertKeyEnvKey   = "PRIVATE_CERT_KEY"
)

// Proxy returns a proxy instance
func Proxy() *goproxy.ProxyHttpServer {
	placeholderImagePath := os.Getenv(placeholderImageEnvKey)
	if placeholderImagePath == "" {
		log.Fatalf("failed to load placeholer image. No %s found", placeholderImageEnvKey)
	}
	placeholder, err := os.Open(placeholderImagePath)
	if err != nil {
		log.Fatalf("failed to load placeholder image: %w", err)
	}
	defer placeholder.Close()

	rawBody, err := io.ReadAll(placeholder)
	if err != nil {
		log.Fatalf("error loading placeholder: %w", err)
	}
	publicCertPath := os.Getenv(publicCertEnvKey)
	privateCertKeyPath := os.Getenv(privateCertKeyEnvKey)
	if publicCertPath == "" || privateCertKeyPath == "" {
		log.Fatalf("both %s and %s environment variables must be set", publicCertPath, privateCertKeyPath)
	}

	cert, err := tls.LoadX509KeyPair(publicCertPath, privateCertKeyPath)
	if err != nil {
		log.Fatalf("error parsing TLS certificate: %w", err)
	}
	proxy := goproxy.NewProxyHttpServer()

	// Go proxies always break websites because they canonicalize HTTP headers by default.
	// While this is a good standard, some other languages and server-client systems don't follow
	// this standard.
	// Therefore, Go proxies break websites.
	// In any go proxy, this auto canonicalization of headers should be disabled.
	proxy.PreventCanonicalization = true
	proxy.KeepDestinationHeaders = true
	proxy.KeepHeader = true
	proxy.AllowHTTP2 = true

	// Allow handling of HTTPS requests by signing them with a man-in-the-middle (MITM) certificate
	customCaMitm := &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&cert)}
	var customAlwaysMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return customCaMitm, host
	}
	proxy.OnRequest().HandleConnect(customAlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if acceptHeader := req.Header.Get("Accept"); strings.Contains(acceptHeader, "image") {
			return req, &http.Response{
				Request:    req,
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(rawBody)),
				Header: http.Header{
					"Content-Type": []string{"image/png"},
				},
			}
		}
		resp, err := ctx.RoundTrip(req)
		if err != nil {
			return req, &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     http.Header{},
				Request:    req,
			}
		}
		return req, resp
	})

	// Enable MITM for HTTPS traffic
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if contentType := resp.Header.Get("Content-Type"); strings.Contains(contentType, "image") {
			log.Println("Intercepting image")
			resp.Body = io.NopCloser(bytes.NewBuffer(rawBody))
		}

		resp.Header.Set("X-Custom-Proxy", "1")
		return resp
	})
	return proxy
}
