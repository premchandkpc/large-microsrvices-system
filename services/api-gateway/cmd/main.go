package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/handler"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/middleware"
	"github.com/premchandkpc/large-microsrvices-system/services/api-gateway/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger := initLogger(cfg)
	defer logger.Sync()

	tp, err := initTracer(cfg)
	if err != nil {
		logger.Fatal("failed to init tracer", zap.Error(err))
	}
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	router := gin.New()
	router.Use(otelgin.Middleware("api-gateway"))
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RateLimiter(cfg))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ready": true})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	svc, err := service.NewServiceRegistry(cfg, logger)
	if err != nil {
		logger.Fatal("failed to create service registry", zap.Error(err))
	}

	h := handler.NewHandler(svc, logger, cfg)

	v1 := router.Group("/api/v1")
	v1.Use(middleware.Authentication(cfg, svc, logger))
	{
		v1.GET("/users/:id", h.GetUser)
		v1.POST("/documents", h.UploadDocument)
		v1.GET("/documents", h.ListDocuments)
		v1.GET("/documents/:id", h.GetDocument)
		v1.POST("/documents/:id/process", h.ProcessDocument)
		v1.GET("/search", h.Search)
		v1.POST("/search/vector", h.VectorSearch)
		v1.GET("/notifications", h.GetNotifications)
		v1.POST("/notifications/read", h.MarkAsRead)
	}

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}

	ws := router.Group("/ws")
	ws.GET("/notifications", h.WebSocketNotifications)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		logger.Info("api-gateway starting", zap.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced shutdown", zap.Error(err))
	}
	logger.Info("server exited")
}

func initLogger(cfg *config.Config) *zap.Logger {
	var logger *zap.Logger
	var err error
	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	return logger
}

func initTracer(cfg *config.Config) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(cfg.OTLPEndpoint),
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("api-gateway"),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
			attribute.String("version", "1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	return tp, nil
}
