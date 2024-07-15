# Используем легковесный базовый образ для финальной стадии
FROM golang:latest as builder
ENV GO111MODULE=on
WORKDIR /build

# Копируем и скачиваем зависимости, чтобы использовать кэш
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и выполняем сборку
COPY . .
RUN make -B all

# Используем легковесный образ для финального контейнера
FROM alpine:latest
WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk --no-cache add bash

# Создаем директорию и копируем файлы из билдера
RUN mkdir -p /app/fixtures
COPY fixtures /app/fixtures
COPY --from=builder /build/bin/api /build/bin/seeder ./
COPY .env.docker ./.env

# Устанавливаем команду по умолчанию (если необходимо)
CMD ["./api"]
