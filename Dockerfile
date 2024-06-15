# Use the official Golang image as the base image
ARG GO_VERSION=1.22.4
FROM golang:${GO_VERSION}-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the project files to the container (main.go and .env are in the root)
COPY . .

# Download dependencies and tidy up
RUN go mod download
RUN go mod tidy

# Build the Go application (main.go is in the root) - specifying target architecture
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/ApiRestFinance -ldflags="-s -w"


# Use a smaller base image for the final image that matches your build architecture
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/ApiRestFinance /app/ApiRestFinance

# Copy additional files if needed (adjust paths if necessary)
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/docs /app/docs


# Command to run your application
ENTRYPOINT ["/app/ApiRestFinance"]
