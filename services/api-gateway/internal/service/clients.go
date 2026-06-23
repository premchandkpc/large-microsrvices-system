package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/gorilla/websocket"
)

type AuthClient interface {
	Login(ctx context.Context, email, password string) (*AuthResponse, error)
	Register(ctx context.Context, email, password, name string) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (string, []string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	Logout(ctx context.Context, userID, refreshToken string)
}

type UserClient interface {
	GetUser(ctx context.Context, userID string) (*UserResponse, error)
}

type DocumentClient interface {
	Upload(ctx context.Context, file io.Reader, header *multipart.FileHeader, userID string) (*DocumentResponse, error)
	List(ctx context.Context, userID, page, pageSize, status string) (interface{}, error)
	Get(ctx context.Context, docID, userID string) (*DocumentResponse, error)
	Process(ctx context.Context, docID, pipeline string) (interface{}, error)
}

type SearchClient interface {
	Search(ctx context.Context, query string, filters []string) (*SearchResponse, error)
	VectorSearch(ctx context.Context, query string, filters []string, topK int) (*SearchResponse, error)
}

type NotificationClient interface {
	GetNotifications(ctx context.Context, userID string) (interface{}, error)
	MarkAsRead(ctx context.Context, userID string, notificationIDs []string) error
	HandleWebSocket(ctx context.Context, conn *websocket.Conn, userID string)
}
