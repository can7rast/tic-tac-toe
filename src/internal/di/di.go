package di

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"log"
	"school21/internal/application"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web"
	"time"
)

var Module = fx.Module("app",
	fx.Provide(
		func() (*datasource.DB, error) {
			dsn := "postgres://postgres:123@localhost:5432/school21?sslmode=disable"
			log.Printf("Подключаемся к БД: %s", dsn)

			// 1. Создаём основной пул — pgxpool (как раньше)
			config, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				return nil, fmt.Errorf("failed to parse config: %w", err)
			}
			config.MaxConns = 20
			config.MinConns = 5
			config.MaxConnLifetime = 10 * time.Minute
			config.MaxConnIdleTime = 5 * time.Minute
			config.HealthCheckPeriod = 1 * time.Minute

			pool, err := pgxpool.NewWithConfig(context.Background(), config)
			if err != nil {
				return nil, fmt.Errorf("failed to create pgxpool: %w", err)
			}

			if err = pool.Ping(context.Background()); err != nil {
				pool.Close()
				return nil, fmt.Errorf("pgxpool ping failed: %w", err)
			}

			// 2. Создаём временный *sql.DB ТОЛЬКО для goose
			sqlDB, err := sql.Open("pgx", dsn)
			if err != nil {
				pool.Close()
				return nil, fmt.Errorf("failed to open sql.DB for migrations: %w", err)
			}

			// Настройки пула для миграций (можно минимальные)
			sqlDB.SetMaxOpenConns(5)
			sqlDB.SetMaxIdleConns(2)

			// 3. Запускаем миграции
			log.Println("Запускаем миграции...")
			goose.SetDialect("postgres")

			if err := goose.Up(sqlDB, "migrations"); err != nil {
				sqlDB.Close()
				pool.Close()
				return nil, fmt.Errorf("failed to apply migrations: %w", err)
			}

			log.Println("Миграции успешно применены ✅")

			// 4. Закрываем временный sql.DB — он нам больше не нужен
			sqlDB.Close()

			// 5. Возвращаем только нашу структуру с pgxpool
			return &datasource.DB{Pool: pool}, nil
		},

		datasource.NewUserRepository,
		datasource.NewGameRepository,

		application.NewAuthService,
		application.NewGameService,
		web.SetupRouter,
	),
	fx.Invoke(func(lc fx.Lifecycle, db *datasource.DB) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Println("Graceful shutdown: начинаем закрытие пула БД...")

				done := make(chan error, 1)

				go func() {
					db.Close()
					done <- nil
				}()

				select {
				case err := <-done:
					log.Println("Пул БД закрыт успешно")
					return err
				case <-ctx.Done():
					log.Println("Таймаут при закрытии пула БД — принудительно выходим (но conn закроются ОС)")
					return ctx.Err()
				}
			},
		})
	}),
	fx.Invoke(func(engine *gin.Engine) {}),
)
