# Use the official Go image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application inside the container
RUN go mod download
RUN go build -o exposer -ldflags "-s -w" main.go
RUN chmod +x exposer
RUN mv /app/exposer /usr/local/bin
RUN mkdir -p ~/.config/uncover
RUN mv /app/provider-config.yaml ~/.config/uncover
RUN sed "s/=/: /g" /app/.env > /app/config.yaml

# Run Command
ENTRYPOINT ["exposer", "-h"]
