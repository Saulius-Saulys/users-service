package repository

import (
	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	logger       *zap.Logger
	gormInstance *gorm.DB
}

func NewUser(logger *zap.Logger, gormInstance *gorm.DB) *User {
	return &User{
		logger:       logger,
		gormInstance: gormInstance,
	}
}

func (u *User) Create(user *dto.CreateUser) (*model.User, error) {
	idUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate uuid")
	}

	userModel := &model.User{
		ID:        idUUID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
	}

	result := u.gormInstance.Create(userModel)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to insert user into database")
	}

	return userModel, nil
}

func (u *User) Update(id string, user *dto.UpdateUser) (*model.User, error) {
	userModel := &model.User{
		ID: id,
	}
	if user.FirstName != nil {
		userModel.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		userModel.LastName = *user.LastName
	}
	if user.Nickname != nil {
		userModel.Nickname = *user.Nickname
	}
	if user.Password != nil {
		userModel.Password = *user.Password
	}
	if user.Email != nil {
		userModel.Email = *user.Email
	}
	if user.Country != nil {
		userModel.Country = *user.Country
	}

	result := u.gormInstance.Updates(userModel)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to update user in database")
	}

	return userModel, nil
}

func (u *User) Delete(id string) error {
	result := u.gormInstance.Delete(&model.User{ID: id})

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete user from database")
	}

	return nil
}

func (u *User) GetByCountry(country dto.Country, page, limit int) ([]model.User, error) {
	var users []model.User

	offset := (page - 1) * limit
	result := u.gormInstance.Where("country = ?", country).Offset(offset).Limit(limit).Find(&users)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to retrieve users from database")
	}

	return users, nil
}
