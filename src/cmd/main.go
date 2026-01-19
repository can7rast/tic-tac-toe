package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/fx"

	"school21/internal/di"
)

func main() {
	app := fx.New(
		di.Module,
		fx.Invoke(runServer),
	)

	app.Run()
}

func runServer(lc fx.Lifecycle, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Сервер запущен на http://localhost:8080")
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("Ошибка запуска сервера: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Останавливаем сервер...")
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	})
}
