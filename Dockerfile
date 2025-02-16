# Start with the official Golang image for building the application
FROM --platform=linux/386 golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Build the application as a static binary for Linux x86
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o /proxy-app .

# Start a new minimal base image
FROM --platform=linux/386 alpine:3.18


# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /proxy-app /app/proxy-app

# Expose the defined port
EXPOSE $PORT

# Run the proxy application
CMD ["/app/proxy-app"]
