# Stage 1: Build the Go application
FROM golang:1.22.7-alpine AS builder

# Argument for the application version
ARG APP_VERSION

# Set the working directory
WORKDIR /app

# Copy Go modules files from the backend directory
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code from the backend/internal directory
COPY ./backend/internal ./internal

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o ./out/service \
    -ldflags="-w -s -extldflags '-static' -X main.VERSION=$APP_VERSION" ./internal

# Stage 2: Create a minimal final image
FROM scratch

# Set the working directory
WORKDIR /app

# Add the binary to the PATH
ENV PATH="${PATH}:/app"

# Copy the built binary from the builder stage
COPY --from=builder /app/out/service ./service

# Expose the necessary ports
EXPOSE 8081
EXPOSE 8080
EXPOSE 8079

# Command to run the service
CMD ["service", "start"]
