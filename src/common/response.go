package common

import (
	"log"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ResultResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewErrorResponse(c *gin.Context, statusCode int, err ErrorResponse) {
	if len(err.Message) == 0 {
		err.Message = "error"
	}
	c.AbortWithStatusJSON(statusCode, err)
}

func NewResultResponse(c *gin.Context, statusCode int, result ResultResponse) {
	log.Print(result.Message)
	if len(result.Message) == 0 {
		result.Message = "successfully"
	}

	c.JSON(statusCode, result)
}
