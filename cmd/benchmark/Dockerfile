# Build stage
FROM golang:1.23.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/benchmark ./cmd/benchmark

# Final stage
FROM alpine:latest
RUN apk add --update --no-cache 7zip
COPY --from=builder /bin/benchmark /app/benchmark
CMD ["./app/benchmark"]