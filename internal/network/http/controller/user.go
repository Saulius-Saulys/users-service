package controller

import (
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/network/http/validation"
	"github.com/Saulius-Saulys/users-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	idParam      = "id"
	countryParam = "country"
	pageQuery    = "page"
	defaultPage  = "1"
	limitQuery   = "limit"
	defaultLimit = "10"
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
	reqDTO := &dto.CreateUser{}
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

func (u *User) Update(ctx *gin.Context) {
	userID := ctx.Param(idParam)
	reqDTO := &dto.UpdateUser{}
	if err, errStatusCode := validation.ValidateJSONBody(ctx, reqDTO); err != nil {
		abortErr := ctx.AbortWithError(errStatusCode, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}

	err := u.userService.Update(userID, reqDTO)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

func (u *User) Delete(ctx *gin.Context) {
	userID := ctx.Param(idParam)
	err := u.userService.Delete(userID)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

func (u *User) GetByCountry(ctx *gin.Context) {
	params, err := u.retrieveGetByCountryParams(ctx)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}
	users, err := u.userService.GetByCountry(dto.Country(params.country), params.page, params.limit)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, users)
}

type getByCountryParams struct {
	country string
	page    int
	limit   int
}

func (u *User) retrieveGetByCountryParams(ctx *gin.Context) (*getByCountryParams, error) {
	country := ctx.Param(countryParam)
	page := ctx.DefaultQuery(pageQuery, defaultPage)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert page to integer")
	}
	limit := ctx.DefaultQuery(limitQuery, defaultLimit)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert limit to integer")
	}

	return &getByCountryParams{
		country: country,
		page:    pageInt,
		limit:   limitInt,
	}, nil
}
