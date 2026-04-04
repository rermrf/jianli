package httpapi

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{Code: 0, Data: data})
}

func Error(c *gin.Context, status, code int, message string) {
	c.JSON(status, ErrorResponse{Code: code, Message: message})
}
