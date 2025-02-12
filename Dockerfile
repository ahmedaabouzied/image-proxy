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

# Set environment variables
ENV PORT=8080
ENV PLACEHOLDER_IMAGE_PATH="/app/placeholder.png"
ENV PUB_CERT_PATH="/app/certs/rootCA.pem"
ENV PRIVATE_CERT_KEY="/app/certs/private_key.pem"

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /proxy-app /app/proxy-app

# Copy certificates directory
COPY --from=builder /app/certs /app/certs

# Copy the placeholder.png file
COPY --from=builder /app/placeholder.png /app/placeholder.png

# Expose the defined port
EXPOSE $PORT

# Run the proxy application
CMD ["/app/proxy-app"]
