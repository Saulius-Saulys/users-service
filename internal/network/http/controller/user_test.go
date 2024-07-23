package controller_test

import (
	"net/http"
	"testing"

	"github.com/Saulius-Saulys/users-service/internal/mocks"
	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/repository"
	"github.com/Saulius-Saulys/users-service/internal/service"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var userController *controller.User
var rabbitMQPublisherMock *mocks.MockRabbitMQPublisher

func TestCreate(t *testing.T) {
	setup(t)
	createUserDTO := dto.CreateUser{
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "test nickname",
		Password:  "test_password",
		Email:     "test@email.com",
		Country:   "LTU",
	}

	t.Run("Create user", func(t *testing.T) {
		router := userRouter(userController)
		r := gofight.New()
		r.POST("/user-service/users/").
			SetJSONInterface(createUserDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusCreated, response.Code)
				userFromDB := model.User{}
				result := gormConn.Where("first_name = ? and email = ?", createUserDTO.FirstName, createUserDTO.Email).Find(&userFromDB)
				assert.Equal(t, int64(1), result.RowsAffected)
				assert.NoError(t, result.Error)

				// compare DTO from request and actual model from DB
				assertDTOToModel(t, createUserDTO, userFromDB)

				gormConn.Delete(userFromDB)
			})
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	userUUID, err := uuid.NewUUID()
	assert.NoError(t, err)
	originalModel := model.User{
		ID:        userUUID.String(),
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "test nickname",
		Password:  "test_password",
		Email:     "test@email.com",
		Country:   "LTU",
	}

	updatedFirstName := "updated name"
	updatedEmail := "updated@email.com"
	updateUserDTO := dto.UpdateUser{
		FirstName: &updatedFirstName,
		Email:     &updatedEmail,
	}

	t.Run("Update user", func(t *testing.T) {
		insertResult := gormConn.Create(&originalModel)
		assert.NoError(t, insertResult.Error)

		router := userRouter(userController)
		r := gofight.New()
		path := "/user-service/users/" + originalModel.ID
		rabbitMQPublisherMock.EXPECT().PublishMessage(mock.Anything, mock.Anything).Return(nil).Times(1)

		r.PUT(path).SetJSONInterface(updateUserDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, response.Code)
				updatedUserFromDB := model.User{}
				selectResult := gormConn.Where("first_name = ? and email = ?", updatedFirstName, updatedEmail).Find(&updatedUserFromDB)

				assert.Equal(t, int64(1), selectResult.RowsAffected)
				assert.NoError(t, selectResult.Error)

				assert.Equal(t, updatedFirstName, updatedUserFromDB.FirstName)
				assert.Equal(t, updatedEmail, updatedUserFromDB.Email)

				gormConn.Delete(originalModel)
			})
	})
}

func assertDTOToModel(t *testing.T, createUserDTO dto.CreateUser, userFromDB model.User) {
	assert.Equal(t, createUserDTO.FirstName, userFromDB.FirstName)
	assert.Equal(t, createUserDTO.LastName, userFromDB.LastName)
	assert.Equal(t, createUserDTO.Nickname, userFromDB.Nickname)
	assert.Equal(t, createUserDTO.Email, userFromDB.Email)
	assert.Equal(t, createUserDTO.Country, userFromDB.Country)
}

func userRouter(userController *controller.User) *gin.Engine {
	router := gin.New()
	userEndpoints := router.Group("/user-service/users/")
	userEndpoints.POST("/", userController.Create)
	userEndpoints.PUT("/:id", userController.Update)

	return router
}

func setup(t *testing.T) {
	logger := zap.NewNop()
	userRepository := repository.NewUser(logger, gormConn)
	rabbitMQPublisherMock = mocks.NewMockRabbitMQPublisher(t)
	userService := service.NewUser(logger, userRepository, rabbitMQPublisherMock)
	userController = controller.NewUser(logger, userService)
}
