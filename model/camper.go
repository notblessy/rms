package model

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type CamperRepository interface {
	FindByID(ctx context.Context, id string) (Camper, error)
	FindAll(ctx context.Context, query CamperQueryInput) ([]Camper, int64, error)
	Create(ctx context.Context, camper CamperInput) error
	Update(ctx context.Context, id string, camper CamperInput) error
	Delete(ctx context.Context, id string) error
}

type Camper struct {
	ID              string          `json:"id"`
	ImageUrl        string          `json:"image_url"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	LicensePlate    string          `json:"license_plate"`
	Year            int             `json:"year"`
	Capacity        int             `json:"capacity"`
	Price           decimal.Decimal `json:"price"`
	Condition       string          `json:"condition"`
	LastMaintenance NullTime        `json:"last_maintenance"`
	Transmission    string          `json:"transmission"`
	FuelType        string          `json:"fuel_type"`
	Drivetrain      string          `json:"drivetrain"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type CamperQueryInput struct {
	Keyword string `query:"keyword"`
	PaginatedRequest
}

type CamperInput struct {
	Camper
	EquipmentIDs []string `json:"equipment_ids"`
}

func (c CamperInput) ToEntity(id string) Camper {
	return Camper{
		ID:              c.ID,
		ImageUrl:        c.ImageUrl,
		Name:            c.Name,
		LicensePlate:    c.LicensePlate,
		Year:            c.Year,
		Capacity:        c.Capacity,
		Price:           c.Price,
		Condition:       c.Condition,
		LastMaintenance: c.LastMaintenance,
	}
}

func (c CamperInput) Equipments() []CamperEquipment {
	equipments := make([]CamperEquipment, len(c.EquipmentIDs))

	for i, id := range c.EquipmentIDs {
		equipments[i] = CamperEquipment{
			CamperID:    c.ID,
			EquipmentID: id,
		}
	}
	return equipments
}
