package repository

import (
	"context"

	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(d *gorm.DB) model.DriverRepository {
	return &driverRepository{
		db: d,
	}
}

func (d *driverRepository) FindByID(ctx context.Context, id string) (model.Driver, error) {
	logger := logrus.WithField("id", id)

	var driver model.Driver
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&driver).Error
	if err != nil {
		logger.Errorf("Error querying driver: %v", err)
		return model.Driver{}, err
	}

	return driver, nil
}
func (d *driverRepository) FindAll(ctx context.Context, query model.DriverQueryInput) ([]model.Driver, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"query": utils.Dump(query),
	})

	var (
		drivers []model.Driver
		total   int64
	)

	qb := d.db.WithContext(ctx).Model(&model.Driver{})

	if query.Keyword != "" {
		qb = qb.Where("name ILIKE ?", "%"+query.Keyword+"%")
	}

	err := qb.Count(&total).Error
	if err != nil {
		logger.Errorf("Error counting drivers: %v", err)
		return nil, 0, err
	}

	err = qb.Scopes(query.Paginated()).Order(query.Sorted()).Find(&drivers).Error
	if err != nil {
		logger.Errorf("Error querying drivers: %v", err)
		return nil, 0, err
	}

	return drivers, total, nil
}

func (d *driverRepository) Create(ctx context.Context, driver model.Driver) error {
	logger := logrus.WithField("driver", utils.Dump(driver))

	err := d.db.WithContext(ctx).Create(&driver).Error
	if err != nil {
		logger.Errorf("Error creating driver: %v", err)
		return err
	}

	return nil
}

func (d *driverRepository) Update(ctx context.Context, id string, driver model.Driver) error {
	logger := logrus.WithFields(logrus.Fields{
		"id":     id,
		"driver": utils.Dump(driver),
	})

	err := d.db.WithContext(ctx).Model(&model.Driver{}).Where("id = ?", id).Updates(driver).Error
	if err != nil {
		logger.Errorf("Error updating driver: %v", err)
		return err
	}

	return nil
}

func (d *driverRepository) Delete(ctx context.Context, id string) error {
	logger := logrus.WithField("id", id)

	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Driver{}).Error
	if err != nil {
		logger.Errorf("Error deleting driver: %v", err)
		return err
	}

	return nil
}
