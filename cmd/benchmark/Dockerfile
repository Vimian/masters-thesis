# Build stage
FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd/benchmark ./cmd/benchmark
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/benchmark ./cmd/benchmark

# Final stage
FROM alpine:3.21.3
RUN apk add --update --no-cache 7zip wine
COPY --from=builder /bin/benchmark /app/benchmark
COPY /cmd/benchmark/assets /app
CMD ["./app/benchmark"]