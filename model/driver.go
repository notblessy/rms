package model

import (
	"context"
	"time"
)

type DriverRepository interface {
	FindByID(ctx context.Context, id string) (Driver, error)
	FindAll(ctx context.Context, query DriverQueryInput) ([]Driver, int64, error)
	Create(ctx context.Context, driver Driver) error
	Update(ctx context.Context, id string, driver Driver) error
	Delete(ctx context.Context, id string) error
}

type Driver struct {
	ID            string    `json:"id"`
	Photo         string    `json:"photo"`
	Name          string    `json:"name"`
	IDNumber      string    `json:"id_number"`
	LicenseNumber string    `json:"license_number"`
	LicenseExpiry string    `json:"license_expiry"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	JoinDate      time.Time `json:"join_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DriverQueryInput struct {
	Keyword string `query:"keyword"`
	PaginatedRequest
}
