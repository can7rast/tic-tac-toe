# Tic-Tac-Toe

> Веб-приложение игры Крестики-Нолики с ИИ на алгоритме Minimax

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=flat&logo=postgresql&logoColor=white)
![Uber FX](https://img.shields.io/badge/Uber_FX-DI-000000?style=flat)

## Описание

Полнофункциональное веб-приложение для игры в крестики-нолики с двумя режимами:
- **Против компьютера** — непобедимый ИИ на алгоритме Minimax с альфа-бета отсечением
- **На двоих** — мультиплеер через общую ссылку

## Демо

```
┌───┬───┬───┐
│ X │ O │ X │
├───┼───┼───┤
│   │ X │   │
├───┼───┼───┤
│ O │   │ X │  ← X побеждает!
└───┴───┴───┘
```

## Технологии

### Backend

| Технология | Описание |
|------------|----------|
| **Go 1.25** | Язык программирования |
| **Gin** | Высокопроизводительный HTTP-фреймворк |
| **Uber FX** | Dependency Injection контейнер |
| **pgx/v5** | PostgreSQL драйвер с пулом соединений |
| **goose** | Миграции базы данных |
| **bcrypt** | Хеширование паролей |
| **UUID** | Генерация уникальных идентификаторов |

### Frontend

| Технология | Описание |
|------------|----------|
| **Vanilla JS** | Без фреймворков, чистый JavaScript |
| **Single-page** | Один HTML файл со всей логикой |
| **CSS Variables** | Современная тёмная тема |
| **Session Storage** | Хранение сессии авторизации |

### База данных

| Технология | Описание |
|------------|----------|
| **PostgreSQL** | Реляционная СУБД |
| **JSONB** | Хранение состояния доски |
| **Индексы** | Оптимизация запросов |

## Архитектура

Проект построен по принципам **Clean Architecture**:

```
src/
├── cmd/
│   └── main.go                 # Точка входа, запуск сервера
├── internal/
│   ├── domain/                 # Бизнес-логика (ядро)
│   │   ├── game.go             # Сущность игры
│   │   ├── board.go            # Игровое поле
│   │   ├── minimax.go          # Алгоритм ИИ
│   │   ├── User.go             # Сущность пользователя
│   │   ├── service.go          # Интерфейс сервиса
        └── GameRepo.go         # Интерфейс репозитория игр
│   │     
│   ├── application/            # Бизнес-сервисы
│   │   ├── auth_service.go     # Аутентификация
│   │   └── game_service.go     # Логика игры
│   ├── infrastructure/         # Внешние зависимости
│   │   └── datasource/
│   │       ├── db.go           # Подключение к БД
│   │       ├── repository.go   # Репозиторий игр
│   │       └── user_repository.go
│   ├── web/                    # HTTP слой
│   │   ├── router.go           # Маршрутизация
│   │   ├── handler/            # Обработчики запросов
│   │   ├── middleware/         # Middleware (auth)
│   │   └── dto/                # Data Transfer Objects
│   └── di/
│       └── di.go               # Dependency Injection модуль
├── frontend/
│   └── index.html              # SPA фронтенд
├── migrations/                 # SQL миграции
│   ├── 20260128215633_create_users.sql
│   └── 20260128220041_create_games.sql
└── pkg/
    └── auth.go                 # Утилиты для паролей
```

## Особенности и фишки

### Алгоритм Minimax с альфа-бета отсечением

ИИ использует классический алгоритм Minimax с оптимизацией:

```go
func minimax(board *[3][3]int, depth int, alpha, beta int, maximizing bool) int {
    // Терминальные состояния
    if score == 10 {
        return score - depth  // Предпочитаем быструю победу
    }
    if score == -10 {
        return score + depth  // Откладываем поражение
    }
    
    // Альфа-бета отсечение для ускорения
    if beta <= alpha {
        return best
    }
    // ...
}
```

**Ключевые особенности:**
- **Непобедимый ИИ** — всегда играет оптимально
- **Предпочтение быстрых побед** — учёт глубины рекурсии
- **Альфа-бета отсечение** — отсекает заведомо плохие ветки

### Uber FX Dependency Injection

Автоматическое внедрение зависимостей:

```go
var Module = fx.Module("app",
    fx.Provide(
        NewDBConnection,
        datasource.NewUserRepository,
        datasource.NewGameRepository,
        application.NewAuthService,
        application.NewGameService,
        web.SetupRouter,
    ),
)
```

### Graceful Shutdown

Корректное завершение работы сервера:

```go
lc.Append(fx.Hook{
    OnStop: func(ctx context.Context) error {
        log.Println("Graceful shutdown: закрытие пула БД...")
        db.Close()
        return nil
    },
})
```

### Автоматические миграции

Миграции применяются при старте приложения:

```go
goose.Up(sqlDB, "migrations")
```

### Пул соединений PostgreSQL

Оптимизированная работа с БД:

```go
config.MaxConns = 20
config.MinConns = 5
config.MaxConnLifetime = 10 * time.Minute
config.HealthCheckPeriod = 1 * time.Minute
```

## API Endpoints

### Аутентификация

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/signup` | Регистрация нового пользователя |
| `POST` | `/login` | Вход (Basic Auth) |

### Игра (требуется авторизация)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/game` | Создать игру на двоих |
| `POST` | `/game?vs=computer` | Создать игру против ИИ |
| `GET` | `/game/:id` | Получить состояние игры |
| `POST` | `/game/:id` | Сделать ход |
| `POST` | `/game/:id/join` | Присоединиться к игре |
| `GET` | `/games` | Список доступных игр |
| `GET` | `/user/:id` | Информация о пользователе |

### Пример запроса

```bash
# Регистрация
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"login": "player1", "password": "secret123"}'

# Создать игру против ИИ
curl -X POST "http://localhost:8080/game?vs=computer" \
  -H "Authorization: Basic $(echo -n 'player1:secret123' | base64)"

# Сделать ход
curl -X POST http://localhost:8080/game/{id} \
  -H "Authorization: Basic $(echo -n 'player1:secret123' | base64)" \
  -H "Content-Type: application/json" \
  -d '{"board": [[1,0,0],[0,0,0],[0,0,0]], "currentPlayer": 1}'
```

## Установка и запуск

### Требования

- Go 1.25+
- PostgreSQL 14+
- Make (опционально)

### 1. Клонирование

```bash
git clone https://github.com/your-repo/tic-tac-toe.git
cd tic-tac-toe/src
```

### 2. Настройка базы данных

```bash
# Создать базу данных
psql -U postgres -c "CREATE DATABASE school21;"

# Или использовать Docker
docker run -d \
  --name postgres-tictactoe \
  -e POSTGRES_PASSWORD=123 \
  -e POSTGRES_DB=school21 \
  -p 5432:5432 \
  postgres:16
```

### 3. Настройка подключения

Отредактируйте строку подключения в `internal/di/di.go`:

```go
dsn := "postgres://postgres:123@localhost:5432/school21?sslmode=disable"
```

### 4. Запуск

```bash
# Установить зависимости
go mod download

# Запустить сервер
go run cmd/main.go
```

Сервер запустится на `http://localhost:8080`

### 5. Открыть в браузере

Перейдите по адресу: http://localhost:8080

## Схема базы данных

### Таблица `users`

| Колонка | Тип | Описание |
|---------|-----|----------|
| `id` | UUID | Первичный ключ |
| `username` | TEXT | Уникальный логин |
| `password_hash` | TEXT | bcrypt хеш пароля |
| `created_at` | TIMESTAMPTZ | Дата создания |

### Таблица `games`

| Колонка | Тип | Описание |
|---------|-----|----------|
| `id` | UUID | Первичный ключ |
| `board` | JSONB | Состояние доски 3x3 |
| `current_player` | INTEGER | 1 = X, 2 = O |
| `player1_id` | UUID | Создатель игры |
| `player2_id` | UUID | Второй игрок (nullable) |
| `state` | INTEGER | Состояние игры |
| `vsai` | BOOLEAN | Игра против ИИ |
| `created_at` | TIMESTAMPTZ | Дата создания |

### Состояния игры (state)

| Значение | Описание |
|----------|----------|
| 0 | Ожидание игроков |
| 1 | Ход игрока X |
| 2 | Ход игрока O |
| 3 | Ничья |
| 4 | Победа X |
| 5 | Победа O |

## Структура фронтенда

Современный тёмный UI с адаптивным дизайном:

- **Авторизация** — вкладки вход/регистрация
- **Выбор режима** — против компьютера или на двоих
- **Игровое поле** — интерактивная доска 3x3
- **Список игр** — доступные игры для присоединения
- **Лог событий** — история действий


## Возможные улучшения

- [ ] JWT токены вместо Basic Auth
- [ ] WebSocket для real-time обновлений
- [ ] Рейтинговая система игроков
- [ ] История игр
- [ ] Уровни сложности ИИ
## Автор

Учебный проект School 21
