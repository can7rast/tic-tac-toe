# Stage 1: Builder — этап сборки, где мы компилируем Go-код
FROM golang:1.25-alpine AS builder
#Базовый образ с Go 1.25 на Alpine (лёгкая Linux-дистрибуция, ~5MB). Если твоя go.mod требует 1.25.2, Docker возьмёт ближайшую версию.

# Устанавливаем пакеты для сборки: git для go get, ca-certificates для HTTPS, tzdata для timezone (если нужно в логах).
RUN apk add --no-cache git ca-certificates tzdata

# Устанавливаем рабочую директорию в /app
WORKDIR /app

# Сначала копируем go.mod и go.sum — это фишка для кэширования: Docker кэширует слои, так что если моды не изменились, зависимости не скачиваются заново (ускоряет билды в 10 раз).
COPY go.mod go.sum ./
RUN go mod download  # Скачиваем все модули (gin, pgx, fx и т.д.).

# Теперь копируем весь исходный код проекта. Твоя структура (internal/, cmd/, pkg/, migrations/, frontend/) скопируется целиком.
COPY . .

# Компилируем бинарник: CGO_ENABLED=0 — отключаем CGO для статической сборки (бинарник работает без libc); GOOS=linux — для Linux в контейнере.
# -o main — имя бинарника (совпадает с cmd/main.go).

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
# Здесь ./cmd/main.go — путь к твоему main.go. Если ошибка, проверь путь.

# Stage 2: Финальный образ — минимальный, без Go и инструментов сборки.
FROM scratch
# Scratch — пустая ОС (0MB), идеально для Go-бинарников.

# Копируем из builder только то, что нужно: бинарник, сертификаты (для HTTPS в pgx или внешних API), timezone.
COPY --from=builder /app/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/

# Копируем фронтенд и миграции — они статические, и твой router.go отдаёт index.html из frontend/.
COPY frontend /frontend
COPY migrations /migrations
# Твои .sql файлы — goose в di.go их использует.

# Говорим Docker, что приложение слушает на порту 8080 (как в main.go).
EXPOSE 8080

# Команда запуска: просто запускаем бинарник.
CMD ["/main"]