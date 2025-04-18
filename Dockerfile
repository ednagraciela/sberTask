# Используем многоэтапную сборку
# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/todo-app .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080
CMD ["./todo-app"]