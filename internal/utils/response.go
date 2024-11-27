package utils

import (
	"github.com/labstack/echo/v4"
)

type JsonResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx echo.Context, statusCode int, message string, data interface{}) error {
	response := JsonResponse{
		Code:    statusCode,
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return ctx.JSON(statusCode, response)
}

func FailResponse(ctx echo.Context, statusCode int, message string) error {
	response := JsonResponse{
		Status:  "fail",
		Message: message,
	}
	return ctx.JSON(statusCode, response)
}
