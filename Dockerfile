# Start from the latest golang base image
FROM golang:latest

# Add maintainer Info
LABEL maintainer="<your-name> <your-email>"

# Set the current working directory in the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o netbox-hosts-configmap-for-coredns .

# Start from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=0 /app/netbox-hosts-configmap-for-coredns .

# Command to run
CMD ["./netbox-hosts-configmap-for-coredns"]
