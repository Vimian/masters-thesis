# Build stage
FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd/cloud-benchmark ./cmd/cloud-benchmark
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/cloud-benchmark ./cmd/cloud-benchmark

# Final stage
FROM alpine:3.21.3
COPY --from=builder /bin/cloud-benchmark /app/cloud-benchmark
RUN mkdir ./downloads
CMD ["./app/cloud-benchmark"]