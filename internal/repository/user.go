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

func (ur *User) Create(user *dto.User) error {
	idUUID, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "failed to generate uuid")
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

	result := ur.gormInstance.Create(userModel)

	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to insert user into database")
	}

	return nil
}
