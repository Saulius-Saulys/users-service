package controller

import (
	"fmt"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/network/http/validation"
	"github.com/Saulius-Saulys/users-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type User struct {
	logger      *zap.Logger
	userService *service.User
}

func NewUser(logger *zap.Logger, userService *service.User) *User {
	return &User{
		logger:      logger,
		userService: userService,
	}
}

func (u *User) Create(ctx *gin.Context) {
	fmt.Println("as cia vidui 1")
	reqDTO := &dto.User{}
	if err, errStatusCode := validation.ValidateJSONBody(ctx, reqDTO); err != nil {
		abortErr := ctx.AbortWithError(errStatusCode, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}
	err := u.userService.Create(reqDTO)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}
	ctx.JSON(http.StatusCreated, nil)
}
