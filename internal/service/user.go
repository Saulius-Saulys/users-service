package service

import (
	"context"
	"time"

	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RabbitMQPublisher interface {
	PublishMessage(ctx context.Context, message *rabbitmq.OutputMessage) error
}

type UserRepository interface {
	Create(user *dto.CreateUser) (*model.User, error)
	Update(id string, user *dto.UpdateUser) (*model.User, error)
	Delete(id string) error
	GetByCountry(country dto.Country, page, limit int) ([]model.User, error)
}

type User struct {
	logger            *zap.Logger
	userRepository    UserRepository
	rabbitMQPublisher RabbitMQPublisher
}

func NewUser(logger *zap.Logger, userRepository UserRepository, rabbitMQPublisher RabbitMQPublisher) *User {
	return &User{
		logger:            logger,
		userRepository:    userRepository,
		rabbitMQPublisher: rabbitMQPublisher,
	}
}

func (u *User) Create(user *dto.CreateUser) (*model.User, error) {
	hashedPassword, err := utils.HashAndSalt([]byte(user.Password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash and salt password")
	}
	user.Password = hashedPassword

	userModel, err := u.userRepository.Create(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user in database")
	}

	return userModel, nil
}

func (u *User) Update(id string, user *dto.UpdateUser) (*model.User, error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if user.Password != nil {
		hashedPassword, err := utils.HashAndSalt([]byte(*user.Password))
		if err != nil {
			return nil, errors.Wrap(err, "failed to hash and salt password")
		}
		user.Password = &hashedPassword
	}

	userModel, err := u.userRepository.Update(id, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	err = u.rabbitMQPublisher.PublishMessage(
		timeoutCtx,
		&rabbitmq.OutputMessage{
			Action: "update",
			User:   *userModel,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to publish message to rabbitmq")
	}

	return userModel, nil
}

func (u *User) Delete(id string) error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := u.userRepository.Delete(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	err = u.rabbitMQPublisher.PublishMessage(
		timeoutCtx,
		&rabbitmq.OutputMessage{
			Action: "deleted",
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to publish message to rabbitmq")
	}
	err = u.userRepository.Delete(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user from database")
	}

	return nil
}

func (u *User) GetByCountry(country dto.Country, page, limit int) ([]model.User, error) {
	userModels, err := u.userRepository.GetByCountry(country, page, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users by country from database")
	}

	return userModels, nil
}
