# ==========================================
# Этап 1: Сборка бинарника (Builder)
# ==========================================
FROM golang:1.27-alpine AS build

WORKDIR /app

# Сначала копируем только файлы зависимостей (для кэширования слоев)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код проекта
COPY . .

# Собираем статический бинарник без CGO (чтобы работал на чистом alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/taskpulse ./cmd/api/main.go


# ==========================================
# Этап 2: Минимальный образ для запуска (Runner)
# ==========================================
FROM alpine:latest

WORKDIR /app

# Устанавливаем ca-certificates, libc и curl для скачивания migrate
RUN apk --no-cache add ca-certificates libc6-compat curl

# Скачиваем golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin

# Копируем только готовый бинарник из первого этапа
COPY --from=build /app/taskpulse .

# Копируем миграции
COPY internal/migrations ./internal/migrations

# Открываем порт приложения
EXPOSE 8080

# Запускаем миграции и приложение
CMD ["sh", "-c", "migrate -path ./internal/migrations -database $DATABASE_URL up && ./taskpulse"]