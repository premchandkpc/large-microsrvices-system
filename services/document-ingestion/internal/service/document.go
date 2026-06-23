package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/kafka"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/model"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/storage"
	"go.uber.org/zap"
)

type DocumentService struct {
	s3       *storage.S3Client
	producer *kafka.Producer
	cfg      *config.Config
	logger   *zap.Logger
}

func NewDocumentService(s3 *storage.S3Client, producer *kafka.Producer, cfg *config.Config, logger *zap.Logger) *DocumentService {
	return &DocumentService{
		s3:       s3,
		producer: producer,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *DocumentService) GetDocument(ctx context.Context, docID, userID string) (*model.Document, error) {
	return nil, fmt.Errorf("not implemented: get from db")
}

func (s *DocumentService) ListDocuments(ctx context.Context, userID, page, pageSize, status string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: list from db")
}

func (s *DocumentService) ProcessDocument(ctx context.Context, docID, pipelineType string) (*model.ProcessResponse, error) {
	jobID := uuid.New().String()

	event := model.ProcessRequest{
		DocumentID:   docID,
		PipelineType: pipelineType,
		Options: map[string]string{
			"job_id": jobID,
		},
	}

	if err := s.producer.Publish(ctx, "document-processing", docID, event); err != nil {
		s.logger.Error("failed to publish process event", zap.Error(err))
		return nil, fmt.Errorf("failed to queue processing: %w", err)
	}

	s.logger.Info("document queued for processing",
		zap.String("doc_id", docID),
		zap.String("pipeline", pipelineType),
		zap.String("job_id", jobID),
	)

	return &model.ProcessResponse{
		JobID:  jobID,
		Status: model.StatusProcessing,
	}, nil
}
