package controller

import (
	"comedians/src/common"
	"comedians/src/core/users/service"
	"comedians/src/core/usersConcerts/model"
	"comedians/src/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := service.GetUsers()

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: users})
}

func UpdateUser(c *gin.Context) {
	var input model.User
	id, _ := c.Get("userId")

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect data"})

		return
	}

	err := service.UpdateUser(id.(uint64), input)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func DeleteUser(c *gin.Context) {
	id, _  := c.Get("userId")

	err := service.DeleteUser(id.(uint64))

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func GetUser(c *gin.Context) {
	id, _ := c.Get("userId")

	user, err := service.GetUser(id.(uint64))

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: user})
}

func AppendFavoriteConcert(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect concert id"})
		return
	}

	err = service.AppendFavoriteConcert(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending favorite concert"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func AppendFavoriteComedian(c *gin.Context) {
	userId, _ := c.Get("userId")
	comedianId, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect comedian id"})
		return
	}

	err = service.AppendFavoriteComedian(userId.(uint64), comedianId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending favorite comedian"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func AppendSubscription(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect concert id"})
		return
	}

	err = service.AppendSubscription(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending subscription"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func DeleteFavoriteConcert(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect concert id"})
		return
	}

	err = service.DeleteFavoriteConcert(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending favorite concert"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func DeleteFavoriteComedian(c *gin.Context) {
	userId, _ := c.Get("userId")
	comedianId, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect comedian id"})
		return
	}

	err = service.DeleteFavoriteComedian(userId.(uint64), comedianId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending favorite comedian"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func DeleteSubscription(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "Incorrect concert id"})
		return
	}

	err = service.DeleteSubscription(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "Error appending subscription"})
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	group := server.Group("users", middleware.AuthMiddleware)

	group.GET("/", GetUsers)

	currentUserGroup := group.Group("/current")

	currentUserGroup.PUT("/", UpdateUser)
	currentUserGroup.DELETE("/", DeleteUser)
	currentUserGroup.GET("/", GetUser)

	currentUserGroup.POST("/favorite-concerts/:concertId", AppendFavoriteConcert)
	currentUserGroup.POST("/subscriptions/:concertId", AppendSubscription)
	currentUserGroup.POST("/favorite-comedians/:userId", AppendFavoriteComedian)

	currentUserGroup.DELETE("/favorite-concerts/:concertId", DeleteFavoriteConcert)
	currentUserGroup.DELETE("/subscriptions/:concertId", DeleteSubscription)
	currentUserGroup.DELETE("/favorite-comedians/:userId", DeleteFavoriteComedian)

	return group
}
