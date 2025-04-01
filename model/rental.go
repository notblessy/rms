package model

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	RentalStatusPending   = "pending"
	RentalStatusConfirmed = "confirmed"
	RentalStatusCancelled = "cancelled"
	RentalStatusCompleted = "completed"
)

type RentalRepository interface {
	FindByID(ctx context.Context, id string) (Rental, error)
	FindAll(ctx context.Context, query RentalQueryInput) ([]Rental, int64, error)
	Create(ctx context.Context, rental RentalInput) error
	Update(ctx context.Context, id string, rental RentalInput) error
}

type Rental struct {
	ID         string          `json:"id"`
	CustomerID string          `json:"customer_id"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    time.Time       `json:"end_date"`
	RentalType string          `json:"rental_type"`
	CamperID   string          `json:"camper_id"`
	DriverID   string          `json:"driver_id"`
	Status     string          `json:"status"`
	GrandTotal decimal.Decimal `json:"grand_total"`
	Discount   decimal.Decimal `json:"discount"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  NullTime        `json:"deleted_at"`
}

type RentalQueryInput struct {
	Keyword string `query:"keyword"`
	PaginatedRequest
}

type RentalInput struct {
	Rental
	EquipmentIDs []string `json:"equipment_ids"`
}

func (r RentalInput) ToEntity(id string) Rental {
	return Rental{
		ID:         id,
		CustomerID: r.CustomerID,
		StartDate:  r.StartDate,
		EndDate:    r.EndDate,
		RentalType: r.RentalType,
		CamperID:   r.CamperID,
		DriverID:   r.DriverID,
		Status:     r.Status,
		GrandTotal: r.GrandTotal,
		Discount:   r.Discount,
	}
}

func (r RentalInput) Equipments() []RentalEquipment {
	var rentalEquipments []RentalEquipment
	for _, equipmentID := range r.EquipmentIDs {
		rentalEquipments = append(rentalEquipments, RentalEquipment{
			RentalID:    r.ID,
			EquipmentID: equipmentID,
		})
	}
	return rentalEquipments
}

type RentalEquipment struct {
	RentalID    string `json:"rental_id"`
	EquipmentID string `json:"equipment_id"`
}
