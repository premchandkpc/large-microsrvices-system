package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/premchandkpc/large-microsrvices-system/services/document-ingestion/internal/service"
	"go.uber.org/zap"
)

func UploadDocument(svc *service.UploadService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
			return
		}
		defer file.Close()

		limitedReader := io.LimitReader(file, svc.Cfg.MaxUploadSize)
		result, err := svc.Upload(c.Request.Context(), limitedReader, header, userID)
		if err != nil {
			logger.Error("upload failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func GetDocument(svc *service.DocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		docID := c.Param("id")
		userID := c.GetString("user_id")

		doc, err := svc.GetDocument(c.Request.Context(), docID, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		c.JSON(http.StatusOK, doc)
	}
}

func ListDocuments(svc *service.DocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("page_size", "20")
		status := c.Query("status")

		docs, err := svc.ListDocuments(c.Request.Context(), userID, page, pageSize, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, docs)
	}
}

func ProcessDocument(svc *service.DocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		docID := c.Param("id")
		var req struct {
			Pipeline string `json:"pipeline" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pipeline type required"})
			return
		}

		resp, err := svc.ProcessDocument(c.Request.Context(), docID, req.Pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, resp)
	}
}
