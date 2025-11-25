FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-service main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-service ./api-service

EXPOSE 8081

CMD ["./api-service"]
