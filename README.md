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
./image-proxy --port=8080 --stock-image="https://example.com/stock.jpg"
```

Available flags:

```
    --port (default: 8080): The port on which the proxy server listens.
    --stock-image: URL of the image that will be served instead of the original.
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

### Configuration

You can modify the proxy behavior by adjusting:

- The stockImageURL variable in the source code (if hardcoded).
- Environment variables (if implemented).
- Command-line flags as shown above.
