# Image Intercept Proxy

## Overview

`Image Intercept Proxy` is a Golang-based HTTP proxy server designed to intercept image requests and return a predefined stock image instead. This tool is useful for blocking sensitive visual content while maintaining the structure of web pages.

## Features

- Intercepts image requests (`.jpg`, `.png`, `.gif`, etc.)
- Returns a predefined stock image instead of the requested one
- Lightweight and efficient, built with Go's `net/http` package
- Configurable proxy settings

## Installation

### Prerequisites

- Go 1.18+ installed on your system
- (Optional) A predefined stock image hosted on a publicly accessible URL

### Clone the Repository

```sh
git clone https://github.com/yourusername/image-intercept-proxy.git
cd image-intercept-proxy
```

### Usage

#### Running the Proxy Server

To start the proxy server:

```
./image-proxy --port=8080"
```

Available flags:

```
    --port (default: 8080): The port on which the proxy server listens.
```

#### Using the Proxy

Set your browser or system to use the proxy:

1. Open your network settings.
2. Configure the HTTP proxy to http://localhost:8080.
3. Load a webpage that contains images—intercepted images will be replaced with the stock image.

### Testing

Run unit tests to ensure everything is working correctly:

```
go test ./...
```

To test the proxy manually:

- Start the proxy server.

- Use curl to request an image through the proxy:
```
curl -x http://localhost:8080 https://example.com/sensitive-image.jpg -o output.jpg
````
-The saved output.jpg should be the predefined stock image.

# Deployment to a remote host

**Summary Report: The Impact of Running a MITM Proxy on a Smaller Machine**

When running a Go-based Man-in-the-Middle (MITM) proxy (such as one built on top of the [goproxy](https://github.com/elazarl/goproxy) library) on a resource-constrained or smaller CPU machine, several performance bottlenecks are likely to arise. Below is an overview of the main reasons why using a small machine can make the proxy slow or unresponsive, especially under moderate or heavy traffic loads.

---

## 1. High CPU Usage from TLS Operations

In an MITM setup, every request and response involves **two** separate TLS connections:
1. **Client → Proxy**  
2. **Proxy → Origin Server**  

This doubling of encryption and decryption amplifies CPU usage. On a machine with fewer CPU cores or lower clock speeds, this cryptographic overhead can quickly max out available CPU resources.

**Why it’s worse on a smaller machine**:
- Fewer cores available to parallelize cryptographic tasks.  
- Lower clock speeds cause operations like key exchanges and data encryption/decryption to take more time per request.  
- As concurrency grows (multiple users or simultaneous requests), CPU usage spikes to 100%, causing sluggish performance.

---

## 2. Certificate Generation Overhead

By design, a MITM proxy must generate a **per-domain certificate** to present to the client whenever it intercepts HTTPS traffic. If the proxy is not caching these certificates (via a “cert store”), it may repeatedly generate new certificates, which is a **CPU-intensive** task. 

**Why it’s worse on a smaller machine**:
- Generating RSA or ECDSA keys is computationally expensive.  
- A small CPU can become saturated if multiple domains or multiple concurrent requests trigger repeated certificate generation.  

---

## 3. Limited Keep-Alive/Connection Reuse Benefits

Even with HTTP/1.1 or HTTP/2 keep-alive, a smaller machine can struggle to maintain numerous open connections without hitting CPU or memory constraints—particularly because **each** connection is being decrypted and re-encrypted. Larger or more modern machines can handle this easily due to more cores, higher clock speeds, and better hardware acceleration for cryptography.

---

## 4. Concurrency & Thread Scheduling

Go is efficient at managing concurrency, but any environment with minimal CPU resources can experience context-switch overhead if the proxy tries to handle many requests at once. With only one or two cores, the CPU frequently switches goroutines, and cryptographic tasks compound this load.

---

## 5. Expected Slowdowns & Workarounds

Given the above constraints, it’s **expected** for a MITM proxy to run slower on a lower-spec machine. Some potential workarounds or optimizations include:

1. **Certificate Caching**: Ensure the proxy is **not** regenerating domain certificates on every request.  
2. **Smaller Key Sizes or ECDSA**: Use 2048-bit RSA or an ECDSA certificate to reduce cryptographic overhead.  
3. **Enable Keep-Alives**: Make sure you’re reusing client and server connections to reduce new TLS handshakes.  
4. **Limit Interception**: Intercept only the traffic or domains you truly need, so not every request is fully MITMed.  
5. **Move to a Bigger Machine**: Ultimately, if throughput needs are high, using a machine with more CPU power is the most straightforward solution.

---

## Conclusion

Running a full MITM proxy on a small, low-spec machine is inherently CPU-intensive due to the additional cryptographic load. **Slow performance and high CPU usage** are normal and expected in such scenarios. For higher throughput or smoother handling of encrypted traffic, a larger machine with more powerful CPU resources is strongly recommended.
