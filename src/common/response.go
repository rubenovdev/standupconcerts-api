package common

import (
	"github.com/gin-gonic/gin"
)

type ResultResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewErrorResponse(c *gin.Context, statusCode int, err error) {
	res := gin.H {
		"message": err.Error(),
	}

	c.AbortWithStatusJSON(statusCode, res)
}

func NewResultResponse(c *gin.Context, statusCode int, result ResultResponse) {
	if len(result.Message) == 0 {
		result.Message = "successfully"
	}

	c.JSON(statusCode, result)
}
