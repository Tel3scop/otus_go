# Собираем в гошке
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта в контейнер
COPY . .

# Устанавливаем необходимые зависимости для сборки
RUN apk add --no-cache build-base

# Собираем приложение
RUN go build -o ./bin/server ./cmd/calendar/scheduler/main.go

# На выходе тонкий образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранный бинарник из стадии сборки
COPY --from=builder /app/bin/server .

# Копируем конфигурационный файл
ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

# Устанавливаем метки
LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

# Запускаем приложение
CMD ["./server", "-config", "/etc/calendar/config.yaml"]