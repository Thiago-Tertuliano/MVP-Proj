package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/consensys/gnark-crypto/field/pool"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/config"
	"github.com/thiago-tertuliano/estudos-platform/internal/infrastructure/persistence/postgres"
	"github.com/thiago-tertuliano/estudos-platform/internal/presentation/http/router"
	"honnef.co/go/tools/config"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.TextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	pool := postgres.NewConnection(cfg)
	defer pool.Close()

	srv := &http.Server{
		Addr: 		  ":" + cfg.AppPort,
		Handler: 	  router.New(cfg, pool),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// sobe o servidor em goroutine para não bloquear o graceful shutdown

	go func ()  {
		slog.Info("servidor iniciado", "porta", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "erro", err)
			os.Exit(1)
		}
	}()

	// aguarda Ctrl+C / SIGTERM e desliga graciosamente
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("desligando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("falha no shutdown", "erro", err)
	}
}