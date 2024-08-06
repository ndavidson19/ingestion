FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rest-ingestion-service ./cmd/rest-ingestion-service

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/rest-ingestion-service .
EXPOSE 8080
CMD ["./rest-ingestion-service"]
