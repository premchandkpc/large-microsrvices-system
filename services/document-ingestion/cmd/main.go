package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/handler"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/kafka"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/service"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/storage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger := initLogger(cfg)

	s3Client, err := storage.NewS3Client(cfg)
	if err != nil {
		logger.Fatal("failed to create s3 client", zap.Error(err))
	}

	producer, err := kafka.NewProducer(cfg.KafkaBrokers, logger)
	if err != nil {
		logger.Fatal("failed to create kafka producer", zap.Error(err))
	}
	defer producer.Close()

	uploadSvc := service.NewUploadService(s3Client, cfg, logger)
	docSvc := service.NewDocumentService(s3Client, producer, cfg, logger)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "document-ingestion",
			"s3":      s3Client.Health(),
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1/documents")
	{
		v1.POST("/upload", handler.UploadDocument(uploadSvc, logger))
		v1.GET("/:id", handler.GetDocument(docSvc, logger))
		v1.GET("", handler.ListDocuments(docSvc, logger))
		v1.POST("/:id/process", handler.ProcessDocument(docSvc, logger))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("document-ingestion starting", zap.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func initLogger(cfg *config.Config) *zap.Logger {
	logger, _ := zap.NewProduction()
	if cfg.Environment == "development" {
		logger, _ = zap.NewDevelopment()
	}
	return logger
}
