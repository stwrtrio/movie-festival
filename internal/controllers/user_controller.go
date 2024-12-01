package controllers

import (
	"net/http"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
	"github.com/stwrtrio/movie-festival/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) Login(ctx echo.Context) error {
	var userRequest models.UserLoginRequest

	cx := ctx.Request().Context()

	if err := ctx.Bind(&userRequest); err != nil {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if err := ctx.Validate(&userRequest); err != nil {
		return utils.FailResponse(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := c.service.Login(cx, userRequest.Username, userRequest.Password)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusUnauthorized, "invalid credentials")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Access granted", map[string]string{"token": token})
}
