package service

import (
	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/repository"
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

func (u *User) Create(user *dto.User) error {
	return u.userRepository.Create(user)
}
