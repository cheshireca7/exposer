# Go image as the base image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy source code to the container
COPY . .

# Install dependencies
RUN apt-get update && apt-get install -y vim
RUN go install github.com/projectdiscovery/uncover/cmd/uncover@latest

# Build the Go application
RUN go mod download
RUN go build -o /usr/local/bin/exposer -ldflags "-s -w" main.go
RUN chmod +x /usr/local/bin/exposer

# Generate config.yaml from .env
RUN sed "s/=/: /g" /app/docker/.env > /app/config.yaml

# Command
CMD ["exposer", "-h"]
