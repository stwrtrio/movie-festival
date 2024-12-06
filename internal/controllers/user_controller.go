package controllers

import (
	"net/http"

	"github.com/stwrtrio/movie-festival/internal/middlewares"
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

// Create godoc
// @Summary User Register
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 200 {object} utils.JsonResponse "Success create user"
// @Failure 400 {object} utils.JsonResponse "Invalid input"
// @Failure 401 {object} utils.JsonResponse "Unauthorized"
// @Router /api/user/register [post]
func (c *UserController) Register(ctx echo.Context) error {
	cx := ctx.Request().Context()
	var req models.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if err := c.service.Register(cx, req); err != nil {
		if err.Error() == "username already exists" {
			return utils.FailResponse(ctx, http.StatusConflict, err.Error())
		}
		return utils.FailResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "user registered successfully", nil)
}

// Login godoc
// @Summary User Login
// @Description Authenticate a user and return a JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} utils.JsonResponse "Access granted, includes JWT token"
// @Failure 400 {object} utils.JsonResponse "Invalid input"
// @Failure 401 {object} utils.JsonResponse "Unauthorized"
// @Router /api/user/login [post]
func (c *UserController) Login(ctx echo.Context) error {
	var userRequest models.LoginRequest

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

// Logout godoc
// @Summary User Logout
// @Description Logout an authenticated user by invalidating token.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.JsonResponse "Success create user"
// @Failure 400 {object} utils.JsonResponse "Invalid input"
// @Failure 401 {object} utils.JsonResponse "Unauthorized"
// @Router /api/user/logout [post]
func (c *UserController) Logout(ctx echo.Context) error {
	cx := ctx.Request().Context()
	// Retrieve user claims from context
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return utils.FailResponse(ctx, http.StatusUnauthorized, "User not authenticated")
	}

	err := c.service.Logout(cx, claims)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "user logged out successfully", nil)
}
