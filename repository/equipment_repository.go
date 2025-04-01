package repository

import (
	"context"

	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(d *gorm.DB) model.EquipmentRepository {
	return &equipmentRepository{
		db: d,
	}
}

func (e *equipmentRepository) FindByID(ctx context.Context, id string) (model.Equipment, error) {
	logger := logrus.WithField("id", id)

	var equipment model.Equipment
	err := e.db.WithContext(ctx).Where("id = ?", id).First(&equipment).Error
	if err != nil {
		logger.Errorf("Error querying equipment: %v", err)
		return model.Equipment{}, err
	}

	return equipment, nil
}

func (e *equipmentRepository) FindAll(ctx context.Context, query model.EquipmentQueryInput) ([]model.Equipment, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"query": utils.Dump(query),
	})

	var (
		equipments []model.Equipment
		total      int64
	)

	qb := e.db.WithContext(ctx).Model(&model.Equipment{})

	if query.Keyword != "" {
		qb = qb.Where("name ILIKE ?", "%"+query.Keyword+"%")
	}

	err := qb.Count(&total).Error
	if err != nil {
		logger.Errorf("Error counting equipments: %v", err)
		return nil, 0, err
	}

	err = qb.Scopes(query.Paginated()).Order(query.Sorted()).Find(&equipments).Error
	if err != nil {
		logger.Errorf("Error querying equipments: %v", err)
		return nil, 0, err
	}

	return equipments, total, nil
}

func (e *equipmentRepository) Create(ctx context.Context, equipment model.Equipment) error {
	logger := logrus.WithField("equipment", utils.Dump(equipment))

	err := e.db.WithContext(ctx).Create(&equipment).Error
	if err != nil {
		logger.Errorf("Error creating equipment: %v", err)
		return err
	}

	return nil
}

func (e *equipmentRepository) Update(ctx context.Context, id string, equipment model.Equipment) error {
	logger := logrus.WithField("id", id)

	err := e.db.WithContext(ctx).Model(&model.Equipment{}).Where("id = ?", id).Updates(equipment).Error
	if err != nil {
		logger.Errorf("Error updating equipment: %v", err)
		return err
	}

	return nil
}

func (e *equipmentRepository) Delete(ctx context.Context, id string) error {
	logger := logrus.WithField("id", id)

	err := e.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Equipment{}).Error
	if err != nil {
		logger.Errorf("Error deleting equipment: %v", err)
		return err
	}

	return nil
}
