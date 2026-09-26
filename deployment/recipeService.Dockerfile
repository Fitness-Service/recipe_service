FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o recipe-service ./cmd

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/recipe-service .
EXPOSE 8082
CMD ["./recipe-service"]