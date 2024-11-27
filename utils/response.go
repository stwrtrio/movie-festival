package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JsonResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx echo.Context, responseCode int, message string, data interface{}) error {
	response := JsonResponse{
		Code:    responseCode,
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return ctx.JSON(http.StatusOK, response)
}

func FailResponse(ctx echo.Context, message string, statusCode int) error {
	response := JsonResponse{
		Status:  "fail",
		Message: message,
	}
	return ctx.JSON(statusCode, response)
}
