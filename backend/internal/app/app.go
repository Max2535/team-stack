package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/team-stack/backend/internal/app/adapters/cache"
	dbad "github.com/example/team-stack/backend/internal/app/adapters/db"
	"github.com/example/team-stack/backend/internal/app/adapters/event"
	"github.com/example/team-stack/backend/internal/app/adapters/jwt"
	"github.com/example/team-stack/backend/internal/app/core/user"
	"github.com/example/team-stack/backend/internal/config"
	"github.com/example/team-stack/backend/internal/db"
	"github.com/example/team-stack/backend/internal/http"
	"github.com/example/team-stack/backend/internal/telemetry"
	"github.com/gofiber/fiber/v2"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := telemetry.NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			log.Warnw("failed to sync logger", "error", err)
		}
	}()

	_, shutdown, err := telemetry.InitTracer(cfg)
	if err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Warnw("failed to shutdown tracer", "error", err)
		}
	}()

	pg, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pg.Close()

	cacheLayer := cache.NewRedis(cfg.RedisAddr)
	eventBus := event.NewKafka(cfg.KafkaBrokers, cfg.KafkaTopic)
	jwtm := jwt.NewJWTManager(cfg)

	userRepo := dbad.NewUserPostgresRepo(pg)
	userSvc := user.NewService(userRepo, jwtm, eventBus)

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	http.RegisterRoutes(app, cfg, log, userSvc, jwtm, cacheLayer)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listen(cfg.Addr())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return app.ShutdownWithContext(ctx)
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	}
}
