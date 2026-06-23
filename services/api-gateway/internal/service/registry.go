package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserResponse struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type DocumentResponse struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type SearchResult struct {
	DocumentID string            `json:"document_id"`
	Score      float64           `json:"score"`
	Fields     map[string]string `json:"fields"`
	Snippet    string            `json:"snippet"`
}

type SearchResponse struct {
	Results   []SearchResult `json:"results"`
	TotalHits int            `json:"total_hits"`
	TookMs    float64        `json:"took_ms"`
}

type NotificationResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type ServiceRegistry struct {
	auth         AuthClient
	user         UserClient
	document     DocumentClient
	search       SearchClient
	notification NotificationClient
}

func NewServiceRegistry(cfg interface{}, logger interface{}) (*ServiceRegistry, error) {
	l, ok := logger.(*zap.Logger)
	if !ok {
		return nil, fmt.Errorf("logger must be *zap.Logger")
	}
	l.Warn("ServiceRegistry initialized without gRPC clients — all service calls will fail")
	l.Warn("Implement gRPC client connections in NewServiceRegistry before production use")
	return &ServiceRegistry{}, fmt.Errorf("gRPC clients not configured: implement client connections in NewServiceRegistry")
}

func (s *ServiceRegistry) errNotConnected(name string) error {
	return fmt.Errorf("%s: gRPC client not initialized", name)
}

func (s *ServiceRegistry) ValidateToken(ctx context.Context, token string) (string, []string, error) {
	if s.auth == nil {
		return "", nil, s.errNotConnected("auth")
	}
	return s.auth.ValidateToken(ctx, token)
}

func (s *ServiceRegistry) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	if s.auth == nil {
		return nil, s.errNotConnected("auth")
	}
	return s.auth.Login(ctx, email, password)
}

func (s *ServiceRegistry) Register(ctx context.Context, email, password, name string) (*AuthResponse, error) {
	if s.auth == nil {
		return nil, s.errNotConnected("auth")
	}
	return s.auth.Register(ctx, email, password, name)
}

func (s *ServiceRegistry) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	if s.auth == nil {
		return nil, s.errNotConnected("auth")
	}
	return s.auth.RefreshToken(ctx, refreshToken)
}

func (s *ServiceRegistry) Logout(ctx context.Context, userID, refreshToken string) {
	if s.auth == nil {
		return
	}
	s.auth.Logout(ctx, userID, refreshToken)
}

func (s *ServiceRegistry) GetUser(ctx context.Context, userID, requesterID string, roles []string) (*UserResponse, error) {
	if s.user == nil {
		return nil, s.errNotConnected("user")
	}
	return s.user.GetUser(ctx, userID)
}

func (s *ServiceRegistry) UploadDocument(ctx context.Context, file io.Reader, header *multipart.FileHeader, userID string) (*DocumentResponse, error) {
	if s.document == nil {
		return nil, s.errNotConnected("document")
	}
	return s.document.Upload(ctx, file, header, userID)
}

func (s *ServiceRegistry) ListDocuments(ctx context.Context, userID, page, pageSize, status string) (interface{}, error) {
	if s.document == nil {
		return nil, s.errNotConnected("document")
	}
	return s.document.List(ctx, userID, page, pageSize, status)
}

func (s *ServiceRegistry) GetDocument(ctx context.Context, docID, userID string) (*DocumentResponse, error) {
	if s.document == nil {
		return nil, s.errNotConnected("document")
	}
	return s.document.Get(ctx, docID, userID)
}

func (s *ServiceRegistry) ProcessDocument(ctx context.Context, docID, pipeline string) (interface{}, error) {
	if s.document == nil {
		return nil, s.errNotConnected("document")
	}
	return s.document.Process(ctx, docID, pipeline)
}

func (s *ServiceRegistry) Search(ctx context.Context, query string, filters []string) (*SearchResponse, error) {
	if s.search == nil {
		return nil, s.errNotConnected("search")
	}
	return s.search.Search(ctx, query, filters)
}

func (s *ServiceRegistry) VectorSearch(ctx context.Context, query string, filters []string, topK int) (*SearchResponse, error) {
	if s.search == nil {
		return nil, s.errNotConnected("search")
	}
	return s.search.VectorSearch(ctx, query, filters, topK)
}

func (s *ServiceRegistry) GetNotifications(ctx context.Context, userID string) (interface{}, error) {
	if s.notification == nil {
		return nil, s.errNotConnected("notification")
	}
	return s.notification.GetNotifications(ctx, userID)
}

func (s *ServiceRegistry) MarkNotificationsRead(ctx context.Context, userID string, notificationIDs []string) error {
	if s.notification == nil {
		return s.errNotConnected("notification")
	}
	return s.notification.MarkAsRead(ctx, userID, notificationIDs)
}

func (s *ServiceRegistry) HandleWebSocket(ctx context.Context, conn *websocket.Conn, userID string) {
	if s.notification == nil {
		return
	}
	s.notification.HandleWebSocket(ctx, conn, userID)
}
