package main_test

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	imgproxy "github.com/ahmedaabouzied/image-proxy"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type handler struct{}

// An HTTP handler to write form data and header into the response body.
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
	// Init our service proxy.
	proxy := imgproxy.Proxy()

	// Create a web server from the proxy.
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	// Create a transport layer and assign the proxy.
	proxyURL, _ := url.Parse(proxyServer.URL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			// Skip HTTPS for now.
			InsecureSkipVerify: true,
		},
		// Set our proxy URL.
		Proxy: http.ProxyURL(proxyURL),
	}
	// Create an HTTP client from the transport URL
	client := &http.Client{Transport: tr}

	httpServer := httptest.NewServer(handler{})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, httpServer.URL+"/proxy", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, resp.Header.Get("X-Image-Proxy"), "1")
}

// A custom transport layer that does not canonicalize HTTP headers.
type customTransport struct {
	proxyAddr string
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Connect to the proxy not the request URL
	addr := req.URL.Host
	if t.proxyAddr != "" {
		addr = t.proxyAddr
	}
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rawRequest := ""
	if t.proxyAddr == "" {
		rawRequest += fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL.RequestURI())
	} else {
		rawRequest += fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL.Scheme+"://"+req.URL.Host+req.URL.RequestURI())
	}
	rawRequest += fmt.Sprintf("HOST: %s\r\n", req.URL.Host)
	for key, values := range req.Header {
		for _, value := range values {
			rawRequest += fmt.Sprintf("%s: %s\r\n", key, value)
		}
	}
	rawRequest += "\r\n"

	_, err = conn.Write([]byte(rawRequest))
	if err != nil {
		return nil, err
	}

	// Read response.
	// Example Response:
	// Status line: HTTP/1.1 200 OK
	// Headers:
	//   non-canonical-header-name: 1
	//   Date: Sat, 08 Feb 2025 13:57:45 GMT
	//   Content-Length: 0
	reader := bufio.NewReader(conn)

	// Read status line
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	statusLine = strings.TrimSpace(statusLine)

	var proto string
	var statusCode int
	var statusText string
	_, err = fmt.Sscanf(statusLine, "%s %d %s", &proto, &statusCode, &statusText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status line %s: %w", statusLine, err)
	}

	headers := map[string][]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "\r\n" { // Empty line
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[parts[0]] = []string{strings.TrimSpace(parts[1])}
		}
	}

	return &http.Response{
		Status:     statusText,
		StatusCode: statusCode,
		Proto:      proto,
		Header:     headers,
		Body:       nil,
		Request:    req,
	}, nil
}

// Go proxies always break websites because they canonilize HTTP headers by default.
// While this is a good standard, some other languages and server-client systems don't follow
// this standard.
// Therefore, Go proxies break websites.
// In any go proxy, this auto canonilization of headers should be disabled.
func TestProxy_HeaderCanonicalizationDisabled(t *testing.T) {
	// Init our service proxy.
	proxy := imgproxy.Proxy()

	// Create a web server from the proxy.
	proxyServer := httptest.NewServer(proxy)
	defer proxyServer.Close()

	// Create a transport layer and assign the proxy and create an HTTP client from it.
	client := &http.Client{Transport: &customTransport{
		proxyAddr: strings.TrimPrefix(proxyServer.URL, "http://"),
	}}

	// Space before header name to skip canonicalizing.
	headerName := "non-canonical-header-name"

	nonCanonicalHeaderHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header()[headerName] = []string{"1"}
		w.WriteHeader(200)
	}

	httpServer := httptest.NewServer(http.HandlerFunc(nonCanonicalHeaderHandler))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, httpServer.URL+"/proxy", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	if err != nil {
		t.Fail()
	}

	for name, v := range resp.Header {
		t.Logf("'%s' : %s", name, v)
	}

	assert.Equal(t, "1", resp.Header.Get("X-Image-Proxy"))
	assert.Equal(t, "1", resp.Header.Get(headerName)) // Get() cananonilizes the given header name
}
