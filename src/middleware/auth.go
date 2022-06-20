package middleware

import (
	authService "comedians/src/core/auth/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func AuthMiddleware(c *gin.Context) {
	log.Print("requewst!!")
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(401, gin.H{
			"message": "No auth header",
		})
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		c.JSON(401, gin.H{
			"message": "Invalid auth header",
		})
		c.Abort()

		return
	}

	// parse token

	claims, err := authService.ParseTokenJWT(headerParts[1])

	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}

	c.Set("userId", claims.Id)
	c.Set("userRoles", claims.Roles)
	c.Next()
}
