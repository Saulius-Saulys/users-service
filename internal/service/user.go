package service

import (
	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/repository"
	"github.com/Saulius-Saulys/users-service/internal/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type User struct {
	logger         *zap.Logger
	userRepository *repository.User
	rabbitMQClient rabbitmq.Client
}

func NewUser(logger *zap.Logger, userRepository *repository.User, rabbitMQClient rabbitmq.Client) *User {
	return &User{
		logger:         logger,
		userRepository: userRepository,
		rabbitMQClient: rabbitMQClient,
	}
}

func (u *User) Create(user *dto.CreateUser) error {
	hashedPassword, err := utils.HashAndSalt([]byte(user.Password))
	if err != nil {
		return errors.Wrap(err, "failed to hash and salt password")
	}
	user.Password = hashedPassword

	return u.userRepository.Create(user)
}

func (u *User) Update(id string, user *dto.UpdateUser) error {
	if user.Password != nil {
		hashedPassword, err := utils.HashAndSalt([]byte(*user.Password))
		if err != nil {
			return errors.Wrap(err, "failed to hash and salt password")
		}
		user.Password = &hashedPassword
	}
	return u.userRepository.Update(id, user)
}

func (u *User) Delete(id string) error {
	return u.userRepository.Delete(id)
}

func (u *User) GetByCountry(country dto.Country, page, limit int) ([]model.User, error) {
	return u.userRepository.GetByCountry(country, page, limit)
}
