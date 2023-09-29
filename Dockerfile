# Go image as the base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy source code to the container
COPY . .

# Install dependencies
RUN apt-get update && apt-get install -y vim

# Build the Go application
RUN go mod download
RUN go build -o /usr/local/bin/exposer -ldflags "-s -w" main.go
RUN chmod +x /usr/local/bin/exposer

# Generate config.yaml from .env
RUN mkdir -p /root/.config/exposer

# Command
CMD ["exposer", "-h"]
