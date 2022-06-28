package controller

import (
	"comedians/src/common"
	usersModel "comedians/src/core/users/model"
	"comedians/src/core/users/service"
	"comedians/src/core/usersConcerts/model"
	usersConcertModel "comedians/src/core/usersConcerts/model"
	"comedians/src/middleware"
	"comedians/src/utils"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func updateUserImage(c *gin.Context) {
	userId, _ := c.Get("userId")

	form, err := c.MultipartForm()

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
	}

	file := form.File["image"][0]
	openedFile, _ := file.Open()

	filepath, err := service.UpdateUserImage(userId.(uint64), openedFile, file.Filename)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"filepath": filepath,
	}})
}

func getUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect user id"))
		return
	}

	user, err := service.GetUser(id)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"user": user,
	}})
}

func updateUserPassword(c *gin.Context) {
	id, _ := c.Get("userId")

	var input usersModel.UpdateUserPasswordDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	log.Print(input)

	err := service.UpdateUserPassword(id.(uint64), input.Password)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getUsers(c *gin.Context) {
	users, err := service.GetUsers()

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	filtersRole := c.Query("role")

	if len(filtersRole) != 0 {
		users = utils.Filter(users, func(elem *model.User) bool {
			for _, role := range elem.Roles {
				if role.Title == filtersRole {
					return true
				}
			}
			return false
		})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"users": users,
	}})
}

func updateCurrentUser(c *gin.Context) {
	var input usersConcertModel.User
	id, _ := c.Get("userId")

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	err := service.UpdateUser(id.(uint64), input)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func deleteCurrentUser(c *gin.Context) {
	id, _ := c.Get("userId")

	err := service.DeleteUser(id.(uint64))

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getCurrentUser(c *gin.Context) {
	id, _ := c.Get("userId")

	user, err := service.GetUser(id.(uint64))

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"user": user,
	}})
}

func appendFavoriteConcert(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.AppendFavoriteConcert(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func appendFavoriteComedian(c *gin.Context) {
	userId, _ := c.Get("userId")
	comedianId, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.AppendFavoriteComedian(userId.(uint64), comedianId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func appendSubscription(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.AppendSubscription(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func deleteFavoriteConcert(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.DeleteFavoriteConcert(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func deleteFavoriteComedian(c *gin.Context) {
	userId, _ := c.Get("userId")
	comedianId, err := strconv.ParseUint(c.Param("userId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.DeleteFavoriteComedian(userId.(uint64), comedianId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func deleteSubscription(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err = service.DeleteSubscription(userId.(uint64), concertId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func like(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userId, _ := c.Get("userId")

	if err := service.Like(id, userId.(uint64)); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func dislike(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userId, _ := c.Get("userId")

	if err := service.Dislike(id, userId.(uint64)); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	server.Static(os.Getenv("PUBLIC_USERS_IMAGES_DIR"), os.Getenv("USERS_IMAGES_DIR"))

	usersGroup := server.Group("users")
	{
		usersGroup.GET("", getUsers)
		usersGroup.GET(":id", getUser)
	}
	usersGroupAuth := usersGroup.Group("", middleware.AuthMiddleware)
	{
		usersGroupAuth.PUT(":id/like", like)
		usersGroupAuth.PUT(":id/dislike", dislike)
	}

	currentUserGroup := usersGroupAuth.Group("current")
	{
		currentUserGroup.PUT("", updateCurrentUser)
		currentUserGroup.DELETE("", deleteCurrentUser)
		currentUserGroup.GET("", getCurrentUser)

		currentUserGroup.PUT("image", updateUserImage)
		currentUserGroup.PUT("password", updateUserPassword)

		currentUserGroup.PUT("favorite-concerts/:concertId", appendFavoriteConcert)
		currentUserGroup.PUT("subscriptions/:concertId", appendSubscription)
		currentUserGroup.PUT("favorite-comedians/:userId", appendFavoriteComedian)

		currentUserGroup.DELETE("favorite-concerts/:concertId", deleteFavoriteConcert)
		currentUserGroup.DELETE("subscriptions/:concertId", deleteSubscription)
		currentUserGroup.DELETE("favorite-comedians/:userId", deleteFavoriteComedian)
	}

	return usersGroup
}
