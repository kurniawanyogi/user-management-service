package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-management-service/common"
	"user-management-service/model"
	"user-management-service/service"
)

type IUserDelivery interface {
	Register(c *gin.Context)
	Detail(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	Login(c *gin.Context)
}

type userDelivery struct {
	common   common.IRegistry
	services service.IRegistry
}

func NewUserDelivery(common common.IRegistry, services service.IRegistry) *userDelivery {
	return &userDelivery{
		common:   common,
		services: services,
	}
}

func (d *userDelivery) Register(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload = model.RegistrationUserRequest{}
	)

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate payload request
	err := d.common.GetValidator().Struct(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "Validation failed",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	err = d.services.GetUserService().Register(ctx, payload)
	if err != nil {
		// handle case error validation, return 400
		if err.Error() == common.ErrEmailAlreadyTaken.Error() || err.Error() == common.ErrUsernameAlreadyTaken.Error() {
			c.JSON(http.StatusBadRequest, model.Response{
				Status:  common.StatusError,
				Message: err.Error(),
			})
			return
		}

		// handle case error server, return 500
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "user registered successfully",
	})
	return
}

func (d *userDelivery) Detail(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		id  = c.Param("id")
	)

	// validate parameter id
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid userId",
			Data:    common.FormatValidationError(err),
		})
		return
	}
	if userId < 1 {
		err := errors.New("invalid userId")
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid userId",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	user, err := d.services.GetUserService().Detail(ctx, userId)
	if err != nil {
		// handle case error validation, return 400
		if err.Error() == common.ErrDataNotFound.Error() {
			c.JSON(http.StatusNotFound, model.Response{
				Status:  common.StatusError,
				Message: err.Error(),
			})
			return
		}

		// handle case error server, return 500
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "User detail fetched successfully",
		Data:    user,
	})
}

func (d *userDelivery) Update(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		id      = c.Param("id")
		payload model.UpdateUserRequest
	)

	// validate parameter id
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid user ID",
			Data:    common.FormatValidationError(err),
		})
		return
	}
	if userId < 1 {
		err := errors.New("invalid userId")
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid userId",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate payload request
	err = d.common.GetValidator().Struct(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "Validation failed",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	payload.Id = userId
	err = d.services.GetUserService().Update(ctx, payload)
	if err != nil {
		// handle case error validation, return 400
		if err.Error() == common.ErrDataNotFound.Error() {
			c.JSON(http.StatusNotFound, model.Response{
				Status:  common.StatusError,
				Message: err.Error(),
			})
			return
		}

		// handle case error server, return 500
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "User updated successfully",
	})
}

func (d *userDelivery) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		id  = c.Param("id")
	)

	// validate parameter id
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid userId",
			Data:    common.FormatValidationError(err),
		})
		return
	}
	if userId < 1 {
		err := errors.New("invalid userId")
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  common.StatusError,
			Message: "invalid userId",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	err = d.services.GetUserService().Delete(ctx, userId)
	if err != nil {
		// handle case error validation, return 400
		if err.Error() == common.ErrDataNotFound.Error() || err.Error() == common.ErrUserAlreadyDeleted.Error() {
			c.JSON(http.StatusNotFound, model.Response{
				Status:  common.StatusError,
				Message: err.Error(),
			})
			return
		}

		// handle case error server, return 500
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "User deleted successfully",
	})
}

func (d *userDelivery) List(c *gin.Context) {
	var ctx = c.Request.Context()

	users, err := d.services.GetUserService().List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "User list fetched successfully",
		Data:    users,
	})
}

func (d *userDelivery) Login(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload = model.LoginRequest{}
	)

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate payload request
	err := d.common.GetValidator().Struct(payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Response{
			Status:  common.StatusError,
			Message: "Validation failed",
			Data:    common.FormatValidationError(err),
		})
		return
	}

	token, user, err := d.services.GetUserService().Login(ctx, payload)
	if err != nil {
		// handle case error validation, return 401
		if err.Error() == common.ErrDataNotFound.Error() ||
			err.Error() == common.ErrUserAlreadyDeleted.Error() ||
			err.Error() == common.ErrInvalidPassword.Error() {

			c.JSON(http.StatusUnauthorized, model.Response{
				Status:  common.StatusError,
				Message: "invalid username or password",
			})
			return
		}

		// handle case error server, return 500
		c.JSON(http.StatusInternalServerError, model.Response{
			Status:  common.StatusError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  common.StatusSuccess,
		Message: "login successfully",
		Data: gin.H{
			"user":  user,
			"token": token,
		},
	})
	return
}
