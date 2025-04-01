package model

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type EquipmentRepository interface {
	FindByID(ctx context.Context, id string) (Equipment, error)
	FindAll(ctx context.Context, query EquipmentQueryInput) ([]Equipment, int64, error)
	Create(ctx context.Context, equipment Equipment) error
	Update(ctx context.Context, id string, equipment Equipment) error
	Delete(ctx context.Context, id string) error
}

type Equipment struct {
	ID          string          `json:"id"`
	ImageURL    string          `json:"image_url"`
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Stock       int             `json:"stock"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	Condition   string          `json:"condition"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type EquipmentQueryInput struct {
	Keyword string `query:"keyword"`
	PaginatedRequest
}

type CamperEquipment struct {
	CamperID    string `json:"camper_id"`
	EquipmentID string `json:"equipment_id"`
}
