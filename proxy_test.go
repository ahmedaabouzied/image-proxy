package main_test

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	imgproxy "github.com/ahmedaabouzied/image-proxy"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var sb strings.Builder

	// Write form data
	if err := req.ParseForm(); err != nil {
		panic(err)
	}
	sb.WriteString("formResult")
	sb.WriteString(": ")
	sb.WriteString(req.Form.Get("result"))
	sb.WriteString(";")

	// Write header data
	for name, values := range req.Header {
		for _, value := range values {
			sb.WriteString("header-" + name)
			sb.WriteString(": ")
			sb.WriteString(value)
			sb.WriteString(";")
		}
	}

	_, _ = io.WriteString(w, sb.String())
}

func TestProxy_SetProxyHeader(t *testing.T) {
	proxy := imgproxy.Proxy()
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	proxyURL, _ := url.Parse(proxyServer.URL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{Transport: tr}

	httpServer := httptest.NewServer(handler{})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, httpServer.URL+"/proxy", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, resp.Header.Get("X-Image-Proxy"), "1")
}
