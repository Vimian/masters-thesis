# Build stage
FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd/cloudbenchmark ./cmd/cloudbenchmark
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/cloudbenchmark ./cmd/cloudbenchmark

# Final stage
FROM alpine:3.21.3
COPY --from=builder /bin/cloudbenchmark /app/cloudbenchmark
RUN mkdir ./downloads
CMD ["./app/cloudbenchmark"]