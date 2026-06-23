package model

import "time"

type Document struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Filename    string            `json:"filename"`
	ContentType string            `json:"content_type"`
	SizeBytes   int64             `json:"size_bytes"`
	S3Key       string            `json:"s3_key"`
	Status      string            `json:"status"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type UploadResult struct {
	Document   Document `json:"document"`
	UploadURL  string   `json:"upload_url,omitempty"`
	DownloadURL string  `json:"download_url,omitempty"`
}

type ProcessRequest struct {
	DocumentID   string            `json:"document_id"`
	PipelineType string            `json:"pipeline_type"`
	Options      map[string]string `json:"options"`
	UserID       string            `json:"user_id"`
}

type ProcessResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusUploaded   = "uploaded"

	EventDocumentUploaded   = "document.uploaded"
	EventDocumentProcessed  = "document.processed"
	EventDocumentFailed     = "document.failed"
)
