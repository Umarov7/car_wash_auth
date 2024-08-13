FROM golang:1.22.5 AS builder

WORKDIR /app

COPY . .
COPY .env .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

EXPOSE 8081
EXPOSE 50051

CMD ["./myapp"]