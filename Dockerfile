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
RUN go build -o /usr/local/bin/exposer -ldflags "-s -w" exposer.go
RUN chmod +x /usr/local/bin/exposer
RUN mkdir -p /root/.config/exposer /root/.config/uncover/

# Generate configuration files
RUN echo c2hvZGFuOiBbXQpjZW5zeXM6IFtdCmZvZmE6IFtdCnF1YWtlOiBbXQpodW50ZXI6IFtdCnpvb21leWU6IFtdCm5ldGxhczogW10KY3JpbWluYWxpcDogW10KcHVibGljd3d3OiBbXQpodW50ZXJob3c6IFtdCg== | base64 -d > ~/.config/uncover/provider-config.yaml

# Command
CMD ["exposer", "-h"]
