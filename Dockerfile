# Use the official Golang image as the base image
FROM golang:1.22.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the project files to the container
COPY . .

# Copy the .env file to the container
COPY ./cmd/.env /app/cmd/.env

# Copy the Swagger docs folder
COPY ./docs /app/docs

# Download dependencies and tidy up
RUN go mod download
RUN go mod tidy

# Build the Go application
RUN GOOS=linux GOARCH=arm64 go build -o apiRestFinance ./cmd/

# Use a smaller base image for the final image (scratch is the smallest possible image)
FROM scratch

# Set the working directory
WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/apiRestFinance .

# Copy the .env file from the builder stage to the final image
COPY --from=builder /app/cmd/.env ./cmd/

# Copy the Swagger docs folder from the builder stage
COPY --from=builder /app/docs ./docs

# Expose the port your application listens on
EXPOSE 8080

# Command to run your application
CMD ["/app/apiRestFinance"]
