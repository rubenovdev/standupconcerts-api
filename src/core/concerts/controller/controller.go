package controller

import (
	"comedians/src/common"
	"comedians/src/core/concerts/comments/controller"
	concertsModel "comedians/src/core/concerts/model"
	"comedians/src/core/concerts/service"
	usersConcertsModel "comedians/src/core/usersConcerts/model"
	"comedians/src/middleware"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	createConcertsPermission = "create_concerts"
)

func createConcert(c *gin.Context) {
	c.MultipartForm()

	form, err := c.MultipartForm()
	log.Print(form)

	if err != nil {
		log.Panic(err)

		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["concert"]

	if len(files) == 0 {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("no file"))
		return
	}

	file := files[0]
	openedFile, _ := file.Open()

	filepath, err := service.UploadConcertFile(openedFile, file.Filename)

	if err != nil {
		log.Panic(err)
		common.NewErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	concert := usersConcertsModel.Concert{}
	concert.Filepath = filepath
	concert.Title = c.Request.Form.Get("title")
	concert.Description = c.Request.Form.Get("description")
	userId, _ := c.Get("userId")

	concert.UserId = userId.(uint64)

	if err := service.CreateConcert(concert); err != nil {
		log.Panic(err)

		service.DeleteConcertFile(filepath)
		common.NewErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	common.NewResultResponse(c, http.StatusCreated, common.ResultResponse{Result: gin.H{
		"filepath": filepath,
	}})
}

func getConcerts(c *gin.Context) {
	excludedIdStr, _ := c.GetQuery("excluded_id")
	comedianIdStr, _ := c.GetQuery("comedian_id")
	yearStr, _ := c.GetQuery("year")

	year, _ := strconv.ParseUint(yearStr, 10, 64)
	comedianId, _ := strconv.ParseUint(comedianIdStr, 10, 64)
	sortBy, _ := c.GetQuery("sort_by")
	excludedId, _ := strconv.ParseUint(excludedIdStr, 10, 64)

	filters := concertsModel.Filters{Year: year, ComedianId: comedianId, SortBy: sortBy, ExcludedId: excludedId}

	concerts, err := service.GetConrerts(filters)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusCreated, common.ResultResponse{Result: gin.H{
		"concerts": concerts,
	}})
}

func deleteConcert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	err := service.DeleteConcert(id)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getConcert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	concert, err := service.GetConcert(id)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"concert": concert,
	}})
}

func updateConcert(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	var concert usersConcertsModel.Concert

	if err := c.ShouldBindJSON(&concert); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect concert data"))
		return
	}

	concert.Id = id

	err := service.UpdateConcert(concert)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func like(c *gin.Context) {
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)
	userId, _ := c.Get("userId")

	if err := service.Like(concertId, userId.(uint64)); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}


func incViews(c *gin.Context) {
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err := service.IncViews(concertId); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func handleConcertId(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, errors.New("incorrect concert id"))

		return
	}
}

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	concertGroup := server.Group("concerts", middleware.AuthMiddleware)

	server.Static(os.Getenv("PUBLIC_CONCERTS_DIR"), os.Getenv("CONCERTS_DIR"))

	concertGroup.POST("", middleware.RoleMiddleware([]string{createConcertsPermission}), createConcert)
	concertGroup.GET("", getConcerts)

	currentConcertGroup := concertGroup.Group(":concertId", handleConcertId)

	currentConcertGroup.DELETE("", deleteConcert)
	currentConcertGroup.PUT("", updateConcert)
	currentConcertGroup.GET("", getConcert)

	currentConcertGroup.PUT("like", like)
	currentConcertGroup.PUT("views/inc", incViews)

	controller.InitRoutes(currentConcertGroup)

	return currentConcertGroup
}
