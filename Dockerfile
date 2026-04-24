# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy the modules needed
COPY proto/ ./proto/
COPY services/fleet/ ./services/fleet/

# Create a service-specific go.work to avoid loading other services
RUN printf "go 1.25.0\n\nuse (\n\t./proto\n\t./services/fleet\n)\n" > go.work

# Build the application
WORKDIR /app/services/fleet
RUN go build -o bin/api cmd/api/main.go cmd/api/module.go

# Final stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/services/fleet/bin/api .
COPY --from=builder /app/services/fleet/.env.template .env

EXPOSE 8081 9090

CMD ["./api"]
