// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nustvondev/otp-service/config"
	"github.com/nustvondev/otp-service/internal/controller/http"

	"github.com/nustvondev/otp-service/pkg/httpserver"
	"github.com/nustvondev/otp-service/pkg/logger"
	"github.com/nustvondev/otp-service/pkg/postgres"
	"github.com/nustvondev/otp-service/pkg/redis"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.MinPoolSize(cfg.PG.PoolMin),
		postgres.MaxConnLifetime(cfg.PG.MaxConnLifetime),
		postgres.MaxConnIdleTime(cfg.PG.MaxConnIdleTime),
		postgres.HealthCheckPeriod(cfg.PG.HealthCheckPeriod),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Redis client
	redisClient, err := redis.New(cfg.Redis.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer redisClient.Close()

	// HTTP Server
	httpServer := httpserver.New(cfg, httpserver.Port(cfg.HTTP.Port))
	http.NewRouter(httpServer.App, cfg, l)

	// Start servers
	// rmqServer.Start()
	// grpcServer.Start()
	httpServer.Start()

	l.Info("app running at port http:%s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("%s", "app - Run - signal: "+s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-grpcServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	// err = grpcServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	// }

	// err = rmqServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	// }
}
