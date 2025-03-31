package repository

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type camperRepository struct {
	db *gorm.DB
}

// NewCamperRepository :nodoc:
func NewCamperRepository(d *gorm.DB) model.CamperRepository {
	return &camperRepository{
		db: d,
	}
}

func (c *camperRepository) FindByID(ctx context.Context, id string) (model.Camper, error) {
	logger := logrus.WithField("id", id)

	var camper model.Camper
	err := c.db.WithContext(ctx).Where("id = ?", id).First(&camper).Error
	if err != nil {
		logger.Errorf("Error querying camper: %v", err)
		return model.Camper{}, err
	}

	return camper, nil
}

func (c *camperRepository) FindAll(ctx context.Context, query model.CamperQueryInput) ([]model.Camper, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"page": query.Page,
		"size": query.Size,
	})

	var campers []model.Camper
	var total int64

	qb := c.db.WithContext(ctx).Model(&model.Camper{})

	if query.Keyword != "" {
		qb = qb.Where("name ILIKE ?", "%"+query.Keyword+"%")
	}

	err := qb.Count(&total).Error
	if err != nil {
		logger.Errorf("Error counting campers: %v", err)
		return nil, 0, err
	}

	err = qb.Scopes(query.Paginated()).Order(query.Sorted()).Find(&campers).Error
	if err != nil {
		logger.Errorf("Error querying campers: %v", err)
		return nil, 0, err
	}

	return campers, total, nil
}

func (c *camperRepository) Create(ctx context.Context, camper model.CamperInput) error {
	logger := logrus.WithField("camper", utils.Dump(camper))

	id, err := gonanoid.New()
	if err != nil {
		logger.Errorf("Error generating ID: %v", err)
		return err
	}

	payload := camper.ToEntity(id)

	tx := c.db.WithContext(ctx).Begin()

	err = tx.Create(&payload).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("Error creating camper: %v", err)
		return err
	}

	if len(camper.EquipmentIDs) > 0 {
		equipments := camper.Equipments()

		err = tx.Create(&equipments).Error
		if err != nil {
			logger.Error(err)
			tx.Rollback()

			return err
		}
	}

	tx.Commit()
	return nil
}

func (c *camperRepository) Update(ctx context.Context, id string, camper model.CamperInput) error {
	logger := logrus.WithField("id", id)

	payload := camper.ToEntity(id)

	tx := c.db.WithContext(ctx).Begin()

	err := tx.Model(&model.Camper{}).Where("id = ?", id).Updates(payload).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("Error updating camper: %v", err)
		return err
	}

	if len(camper.EquipmentIDs) > 0 {
		equipments := camper.Equipments()

		err = tx.Where("camper_id = ?", id).Delete(&model.CamperEquipment{}).Error
		if err != nil {
			logger.Error(err)
			tx.Rollback()

			return err
		}

		err = tx.Create(&equipments).Error
		if err != nil {
			logger.Error(err)
			tx.Rollback()

			return err
		}
	}

	tx.Commit()
	return nil
}

func (c *camperRepository) Delete(ctx context.Context, id string) error {
	logger := logrus.WithField("id", id)

	err := c.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Camper{}).Error
	if err != nil {
		logger.Errorf("Error deleting camper: %v", err)
		return err
	}

	return nil
}
