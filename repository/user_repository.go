package repository

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository :nodoc:
func NewUserRepository(d *gorm.DB) model.UserRepository {
	return &userRepository{
		db: d,
	}
}

func (u *userRepository) Authenticate(ctx context.Context, code, requestOrigin string) (model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"code":           code,
		"request_origin": requestOrigin,
	})

	auth, err := u.verifyToken(ctx, code, requestOrigin)
	if err != nil {
		logger.Errorf("Error verifying token: %v", err)
		return model.User{}, err
	}

	id, err := gonanoid.New()
	if err != nil {
		logger.Errorf("Error generating id: %v", err)
		return model.User{}, err
	}

	var authUser model.User
	err = u.db.Where("email = ?", auth.Email).First(&authUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("Error querying user: %v", err)
		return model.User{}, err
	}

	if err == gorm.ErrRecordNotFound {
		authUser = model.User{
			ID:      id,
			Name:    auth.Name,
			Role:    "USER",
			Email:   auth.Email,
			Picture: auth.Picture,
		}

		err = u.db.Create(&authUser).Error
		if err != nil {
			logger.Errorf("Error creating user: %v", err)
			return model.User{}, err
		}
	}

	return authUser, nil
}

func (u *userRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	logger := logrus.WithField("id", id)

	var user model.User
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		logger.Errorf("Error querying user: %v", err)
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) PatchUser(ctx context.Context, id string, user model.User) error {
	logger := logrus.WithField("id", id)

	updatedFields := map[string]interface{}{}

	if user.Name != "" {
		updatedFields["name"] = user.Name
	}

	if user.Phone != "" {
		updatedFields["phone"] = user.Phone
	}

	if user.Address != "" {
		updatedFields["address"] = user.Address
	}

	if user.IDNumber != "" {
		var existingUser model.User

		err := u.db.Where("id_number = ?", user.IDNumber).First(&existingUser).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Errorf("Error querying user: %v", err)
			return err
		}

		if existingUser.ID != id {
			logger.Errorf("ID number already exists for another user: %v", err)
			return model.ErrDuplicateIDNumber
		}
	}

	err := u.db.Model(&user).Where("id = ?", id).Updates(user).Error
	if err != nil {
		logger.Errorf("Error updating user: %v", err)
		return err
	}

	return nil
}

func (u *userRepository) FindAll(ctx context.Context, query model.UserQueryInput) ([]model.User, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"query": utils.Dump(query),
	})

	var (
		users []model.User
		total int64
	)

	qb := u.db.WithContext(ctx).Model(&model.User{})

	if query.Keyword != "" {
		qb = qb.Where("name LIKE ? OR email ILIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.Role != "" {
		qb = qb.Where("role = ?", query.Role)
	}

	if err := qb.Count(&total).Error; err != nil {
		logger.Errorf("Error counting users: %v", err)
		return nil, 0, err
	}

	if err := qb.Scopes(query.Paginated()).Order(query.Sorted()).Find(&users).Error; err != nil {
		logger.Errorf("Error querying users: %v", err)
		return nil, 0, err
	}

	return users, total, nil
}
