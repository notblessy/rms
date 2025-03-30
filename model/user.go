package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	Authenticate(ctx context.Context, code, requestOrigin string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	PatchUser(ctx context.Context, id string, user User) error

	FindAll(ctx context.Context, query UserQueryInput) ([]User, int64, error)
}

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Picture   string         `json:"picture"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	IDNumber  string         `json:"id_number"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type UserQueryInput struct {
	Keyword string `query:"keyword"`
	Role    string `query:"role"`
	PaginatedRequest
}

type Auth struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type AuthRequest struct {
	Code          string `json:"code"`
	RequestOrigin string `json:"request_origin"`
}

type ChangeUsernameRequest struct {
	Username string `json:"username"`
}

type GoogleAuthInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
