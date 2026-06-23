package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/gorilla/websocket"
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
	return &ServiceRegistry{}, nil
}

func (s *ServiceRegistry) ValidateToken(ctx context.Context, token string) (string, []string, error) {
	return s.auth.ValidateToken(ctx, token)
}

func (s *ServiceRegistry) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	return s.auth.Login(ctx, email, password)
}

func (s *ServiceRegistry) Register(ctx context.Context, email, password, name string) (*AuthResponse, error) {
	return s.auth.Register(ctx, email, password, name)
}

func (s *ServiceRegistry) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	return s.auth.RefreshToken(ctx, refreshToken)
}

func (s *ServiceRegistry) Logout(ctx context.Context, userID, refreshToken string) {
	s.auth.Logout(ctx, userID, refreshToken)
}

func (s *ServiceRegistry) GetUser(ctx context.Context, userID, requesterID string, roles []string) (*UserResponse, error) {
	return s.user.GetUser(ctx, userID)
}

func (s *ServiceRegistry) UploadDocument(ctx context.Context, file io.Reader, header *multipart.FileHeader, userID string) (*DocumentResponse, error) {
	return s.document.Upload(ctx, file, header, userID)
}

func (s *ServiceRegistry) ListDocuments(ctx context.Context, userID, page, pageSize, status string) (interface{}, error) {
	return s.document.List(ctx, userID, page, pageSize, status)
}

func (s *ServiceRegistry) GetDocument(ctx context.Context, docID, userID string) (*DocumentResponse, error) {
	return s.document.Get(ctx, docID, userID)
}

func (s *ServiceRegistry) ProcessDocument(ctx context.Context, docID, pipeline string) (interface{}, error) {
	return s.document.Process(ctx, docID, pipeline)
}

func (s *ServiceRegistry) Search(ctx context.Context, query string, filters []string) (*SearchResponse, error) {
	return s.search.Search(ctx, query, filters)
}

func (s *ServiceRegistry) VectorSearch(ctx context.Context, query string, filters []string, topK int) (*SearchResponse, error) {
	return s.search.VectorSearch(ctx, query, filters, topK)
}

func (s *ServiceRegistry) GetNotifications(ctx context.Context, userID string) (interface{}, error) {
	return s.notification.GetNotifications(ctx, userID)
}

func (s *ServiceRegistry) MarkNotificationsRead(ctx context.Context, userID string, notificationIDs []string) error {
	return s.notification.MarkAsRead(ctx, userID, notificationIDs)
}

func (s *ServiceRegistry) HandleWebSocket(ctx context.Context, conn *websocket.Conn, userID string) {
	s.notification.HandleWebSocket(ctx, conn, userID)
}
