package di

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"school21/internal/application"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web"
)

var Module = fx.Module("app",
	fx.Provide(
		func() (*datasource.DB, error) {
			dst := "postgres://postgres:123@localhost:5432/school21?sslmode=disable"
			log.Printf("connecting to Database")
			return datasource.NewDB(context.Background(), dst)
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
