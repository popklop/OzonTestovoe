FROM golang:1.23.4-alpine AS builder
LABEL authors="Abdul_Kamaz"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY config ./config
EXPOSE 8080

CMD ["./server", "--config=config/config.yaml"]
