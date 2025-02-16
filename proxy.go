package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"go.opentelemetry.io/otel"
)

const (
	placeholderImageEnvKey = "PLACEHOLDER_IMAGE_PATH"
	publicCertEnvKey       = "PUB_CERT_PATH"
	privateCertKeyEnvKey   = "PRIVATE_CERT_KEY"
)

type Interceptor struct {
	proxy            *goproxy.ProxyHttpServer
	cert             tls.Certificate
	caCertPool       *x509.CertPool
	client           *http.Client
	transport        *http.Transport
	placeholderImage []byte

	// Interceptors
    matchReqFunc  goproxy.ReqConditionFunc
	interceptReqFunc  func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)
	interceptRespFunc func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
	mitmFunc          goproxy.FuncHttpsHandler
}

func NewInterceptor() *Interceptor {
	i := &Interceptor{}
	i.loadPlaceholderImage()
	i.loadCert()
	i.loadTransport()
	i.loadClient()
	i.loadMITMFunc()
    i.loadMatchReqFunc()
	i.loadInterceptReqFunc()
	i.loadInterceptRespFunc()
	i.loadProxy()
	return i
}

func (i *Interceptor) loadPlaceholderImage() {
	placeholderImagePath := os.Getenv(placeholderImageEnvKey)
	if placeholderImagePath == "" {
		log.Fatalf("failed to load placeholer image. No %s found", placeholderImageEnvKey)
	}
	placeholder, err := os.Open(placeholderImagePath)
	if err != nil {
		log.Fatalf("failed to load placeholder image: %s", err)
	}
	defer placeholder.Close()

	rawBody, err := io.ReadAll(placeholder)
	if err != nil {
		log.Fatalf("error loading placeholder: %s", err)
	}
	i.placeholderImage = rawBody
}

func (i *Interceptor) loadCert() {
	publicCertPath := os.Getenv(publicCertEnvKey)
	privateCertKeyPath := os.Getenv(privateCertKeyEnvKey)
	if publicCertPath == "" || privateCertKeyPath == "" {
		log.Fatalf("both %s and %s environment variables must be set", publicCertPath, privateCertKeyPath)
	}

	cert, err := tls.LoadX509KeyPair(publicCertPath, privateCertKeyPath)
	if err != nil {
		log.Fatalf("error parsing TLS certificate: %s", err)
	}
	caCert, err := os.ReadFile(publicCertPath)
	if err != nil {
		log.Fatalf("error reading ca cert: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	i.cert = cert
	i.caCertPool = caCertPool
}

func (i *Interceptor) loadTransport() {
	transport := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false, // Allows connection reuse
		TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
		},
	}
	i.transport = transport
}

func (i *Interceptor) loadProxy() {
	proxy := goproxy.NewProxyHttpServer()
    proxy.Tr = i.transport

	// Go proxies always break websites because they canonicalize HTTP headers by default.
	// While this is a good standard, some other languages and server-client systems don't follow
	// this standard.
	// Therefore, Go proxies break websites.
	// In any go proxy, this auto canonicalization of headers should be disabled.
	proxy.PreventCanonicalization = true
	proxy.KeepDestinationHeaders = true
	proxy.KeepHeader = true
	proxy.AllowHTTP2 = true

	// Man in the middle function to intercept HTTP requests
    if i.mitmFunc != nil {
	    proxy.OnRequest().HandleConnect(i.mitmFunc)
    }

	// Intercept requests (after HTTPS handshake) to return our placeholder image
	proxy.OnRequest(i.matchReqFunc).DoFunc(i.interceptReqFunc)

	// Intercept responses to return our placeholder image

	proxy.OnResponse().DoFunc(i.interceptRespFunc)
	i.proxy = proxy
}

func (i *Interceptor) loadClient() {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	i.client = client
}

func (i *Interceptor) loadMatchReqFunc() {
    i.matchReqFunc = goproxy.ReqConditionFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) bool {
		return strings.Contains(r.Header.Get("Accept"), "image")	
    })
}

func (i *Interceptor) loadMITMFunc() {
	tlsConfigFunc := goproxy.TLSConfigFromCA(&i.cert)
	customCaMitm := &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: tlsConfigFunc}
	i.mitmFunc = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return customCaMitm, host
	}
}

func (i *Interceptor) loadInterceptReqFunc() {
	i.interceptReqFunc = func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, &http.Response{
			Request:    req,
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(i.placeholderImage)),
			Header: http.Header{
				"Content-Type": []string{"image/png"},
			},
		}
	}
}

func (i *Interceptor) loadInterceptRespFunc() {
	i.interceptRespFunc = func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		start := time.Now()
		tracer := otel.Tracer("proxy")
		_, span := tracer.Start(ctx.Req.Context(), "request")
		defer span.End()
		if contentType := resp.Header.Get("Content-Type"); strings.Contains(contentType, "image") {
			_, span := tracer.Start(ctx.Req.Context(), "intercept")
			resp.Body = io.NopCloser(bytes.NewBuffer(i.placeholderImage))
			span.End()
		}

		resp.Header.Set("X-Custom-Proxy", "1")
		log.Printf("Response intercepted. Duration %s", time.Since(start))
		return resp
	}
}

// Implement http.Handler interface with underlying proxy handler implementation
func (i *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.proxy.ServeHTTP(w, r)
}
