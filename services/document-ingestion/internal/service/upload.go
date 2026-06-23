package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/config"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/kafka"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/model"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/storage"
	"go.uber.org/zap"
)

type UploadService struct {
	s3     *storage.S3Client
	cfg    *config.Config
	logger *zap.Logger
}

func NewUploadService(s3 *storage.S3Client, cfg *config.Config, logger *zap.Logger) *UploadService {
	return &UploadService{s3: s3, cfg: cfg, logger: logger}
}

func (s *UploadService) Upload(ctx context.Context, file io.Reader, header *multipart.FileHeader, userID string) (*model.UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if err := s.validateFile(header); err != nil {
		return nil, err
	}

	docID := uuid.New().String()
	s3Key := fmt.Sprintf("%s/%s/%s%s", userID, time.Now().Format("2006/01/02"), docID, ext)

	result, err := s.s3.Upload(ctx, s3Key, file)
	if err != nil {
		return nil, fmt.Errorf("uploading to s3: %w", err)
	}

	doc := model.Document{
		ID:          docID,
		UserID:      userID,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		SizeBytes:   header.Size,
		S3Key:       s3Key,
		Status:      model.StatusUploaded,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	presignedURL, _ := s.s3.GetPresignedURL(ctx, s3Key, 24*time.Hour)

	return &model.UploadResult{
		Document:   doc,
		UploadURL:  result.Location,
		DownloadURL: presignedURL,
	}, nil
}

func (s *UploadService) validateFile(header *multipart.FileHeader) error {
	if header.Size > s.cfg.MaxUploadSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", header.Size, s.cfg.MaxUploadSize)
	}

	contentType := header.Header.Get("Content-Type")
	for _, allowed := range s.cfg.AllowedTypes {
		if contentType == allowed {
			return nil
		}
	}
	return fmt.Errorf("file type not allowed: %s", contentType)
}
