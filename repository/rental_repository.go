package repository

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type rentalRepository struct {
	db *gorm.DB
}

// NewRentalRepository :nodoc:
func NewRentalRepository(d *gorm.DB) model.RentalRepository {
	return &rentalRepository{
		db: d,
	}
}

// type RentalRepository interface {
// 	FindByID(ctx context.Context, id string) (Rental, error)
// 	FindAll(ctx context.Context, query RentalQueryInput) ([]Rental, int64, error)
// 	Create(ctx context.Context, rental RentalInput) error
// 	Update(ctx context.Context, id string, rental RentalInput) error
// 	Delete(ctx context.Context, id string) error
// }

func (r *rentalRepository) FindByID(ctx context.Context, id string) (model.Rental, error) {
	logger := logrus.WithField("id", id)

	var rental model.Rental
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rental).Error
	if err != nil {
		logger.Errorf("Error querying rental: %v", err)
		return model.Rental{}, err
	}

	return rental, nil
}

func (r *rentalRepository) FindAll(ctx context.Context, query model.RentalQueryInput) ([]model.Rental, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"query": query,
	})

	var (
		rentals []model.Rental
		total   int64
	)

	qb := r.db.WithContext(ctx).Model(&model.Rental{})

	if query.Keyword != "" {
		qb = qb.Where("name ILIKE ?", "%"+query.Keyword+"%")
	}

	err := qb.Count(&total).Error
	if err != nil {
		logger.Errorf("Error counting rentals: %v", err)
		return nil, 0, err
	}

	err = qb.Scopes(query.Paginated()).Order(query.Sorted()).Find(&rentals).Error
	if err != nil {
		logger.Errorf("Error querying rentals: %v", err)
		return nil, 0, err
	}

	return rentals, total, nil
}

func (r *rentalRepository) Create(ctx context.Context, rental model.RentalInput) error {
	logger := logrus.WithField("rental", utils.Dump(rental))

	var activeRentals []model.Rental

	err := r.db.WithContext(ctx).Where("status <> ? AND camper_id = ?", model.RentalStatusCancelled, rental.CamperID).First(&activeRentals).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("Error querying active rentals: %v", err)
		return err
	}

	id, err := gonanoid.New()
	if err != nil {
		logger.Errorf("Error generating ID: %v", err)
		return err
	}

	rentalPayload := rental.ToEntity(id)

	tx := r.db.WithContext(ctx).Begin()

	err = tx.Create(&rentalPayload).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("Error creating rental: %v", err)
		return err
	}

	if len(rental.EquipmentIDs) > 0 {
		equipments := rental.Equipments()
		err := tx.Create(&equipments).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("Error creating rental equipment: %v", err)
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *rentalRepository) Update(ctx context.Context, id string, rental model.RentalInput) error {
	logger := logrus.WithFields(logrus.Fields{
		"id":     id,
		"rental": utils.Dump(rental),
	})

	var existingRental model.Rental
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&existingRental).Error
	if err != nil {
		logger.Errorf("Error querying rental: %v", err)
		return err
	}

	if existingRental.Status != model.RentalStatusPending {
		logger.Errorf("Cannot update cancelled rental: %v", err)
		return model.ErrRentalCancelled
	}

	rentalPayload := rental.ToEntity(id)

	tx := r.db.WithContext(ctx).Begin()

	err = tx.Model(&existingRental).Updates(rentalPayload).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("Error updating rental: %v", err)
		return err
	}

	if len(rental.EquipmentIDs) > 0 {
		err := tx.Where("rental_id = ?", id).Delete(&model.RentalEquipment{}).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("Error deleting rental equipment: %v", err)
			return err
		}

		equipments := rental.Equipments()

		err = tx.Create(&equipments).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf("Error creating rental equipment: %v", err)
			return err
		}
	}

	tx.Commit()
	return nil
}
