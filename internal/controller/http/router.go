package http

import (
	"net/http"
	"time"

	"github.com/nustvondev/otp-service/config"
	"github.com/nustvondev/otp-service/internal/controller/http/middleware"

	"github.com/nustvondev/otp-service/pkg/logger"
	"github.com/nustvondev/otp-service/pkg/profiling"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(app *gin.Engine, cfg *config.Config, l logger.Interface) {
	// Initialize profiler
	profiler := profiling.NewProfiler(l.Zerolog(), cfg.Profiling.Enabled, cfg.Profiling.Path)

	// Middleware
	app.Use(middleware.Logger(l))
	// app.Use(middleware.Recovery(l))

	// Add profiling middleware
	app.Use(middleware.ProfilingMiddleware(profiler, l.Zerolog()))
	app.Use(middleware.ProfilingContextMiddleware(profiler))

	// Add timeout middleware - affects all routes
	timeoutConfig := middleware.TimeoutConfig{
		Timeout: cfg.HTTP.ApiTimeout, // Request-level timeout (shorter than WriteTimeout)
		TimeoutResponse: gin.H{
			"error":     "Request timeout",
			"message":   "The server took too long to process your request",
			"timeout":   cfg.HTTP.ApiTimeout,
			"timestamp": time.Now().Format(time.RFC3339),
		},
		SkipPaths: []string{"/healthz", "/metrics", "/swagger", "/debug"},
	}
	app.Use(middleware.TimeoutMiddleware(timeoutConfig))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		registerAt := func() string {
			if cfg.Metrics.Path != "" {
				return cfg.Metrics.Path
			}
			return "/metrics"
		}

		app.GET(registerAt(), gin.WrapH(promhttp.Handler()))
	}

	// Profiling routes
	if cfg.Profiling.Enabled {
		profiler.SetupRoutes(app)
	}

	// K8s probe
	app.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

}
