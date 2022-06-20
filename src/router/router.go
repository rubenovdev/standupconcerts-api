package router

import (
	auth "comedians/src/core/auth/controller"
	concerts "comedians/src/core/concerts/controller"
	users "comedians/src/core/users/controller"
	"comedians/src/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	auth.InitGroup(router)
	users.InitGroup(router)
	concerts.InitGroup(router)
	
	return router
}
