package router

import (
	auth "comedians/src/core/auth/controller"
	users "comedians/src/core/users/controller"
	concerts "comedians/src/core/concerts/controller"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.New()

	auth.InitGroup(router)
	users.InitGroup(router)
	concerts.InitGroup(router)
	
	return router
}
