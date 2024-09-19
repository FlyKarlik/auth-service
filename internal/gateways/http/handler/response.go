package handler

import "github.com/gin-gonic/gin"

type goodResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errResponse(c *gin.Context, status string, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Status: status, Code: statusCode, Message: err.Error()})
}

func successResponse(c *gin.Context, status string, statusCode int, data interface{}) {
	c.JSON(statusCode, goodResponse{Status: status, Code: statusCode, Data: data})
}
