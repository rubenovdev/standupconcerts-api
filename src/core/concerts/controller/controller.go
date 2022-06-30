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

	if err != nil {
		log.Panic(err)

		common.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["concert"]
	youtubeLink := c.Request.Form.Get("youtubeVideoLink")

	var (
		filepathVideo string
		filepathFrame string
	)

	if len(files) != 0 {
		file := files[0]
		openedFile, _ := file.Open()

		filepathVideo, filepathFrame, err = service.UploadConcertFile(openedFile, file.Filename)

		if err != nil {
			log.Panic(err)
			common.NewErrorResponse(c, http.StatusInternalServerError, err)

			return
		}
	} else if len(youtubeLink) != 0 {
		filepathVideo, filepathFrame, err = service.DownloadVideoFromYoutube(youtubeLink)

		if err != nil {
			log.Panic(err)
			common.NewErrorResponse(c, http.StatusInternalServerError, err)

			return
		}
	} else {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("need a youtube link or file"))
		return
	}

	concert := usersConcertsModel.Concert{}
	concert.VideoSrc = filepathVideo
	concert.FrameSrc = filepathFrame
	concert.Title = c.Request.Form.Get("title")
	concert.Description = c.Request.Form.Get("description")
	userId, _ := c.Get("userId")

	concert.UserId = userId.(uint64)

	if err := service.CreateConcert(concert); err != nil {
		log.Panic(err)

		service.DeleteConcertVideo(filepathVideo)
		service.DeleteConcertFrame(filepathVideo)

		common.NewErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	common.NewResultResponse(c, http.StatusCreated, common.ResultResponse{Result: gin.H{
		"videoSrc": filepathVideo,
		"frameSrc": filepathFrame,
	}})
}

func getConcerts(c *gin.Context) {
	excludedIdStr, _ := c.GetQuery("excluded_id")
	comedianIdStr, _ := c.GetQuery("comedian_id")
	yearStr, _ := c.GetQuery("year")
	limitStr, _ := c.GetQuery("limit")

	year, _ := strconv.ParseUint(yearStr, 10, 64)
	comedianId, _ := strconv.ParseUint(comedianIdStr, 10, 64)
	sortBy, _ := c.GetQuery("sort_by")
	excludedId, _ := strconv.ParseUint(excludedIdStr, 10, 64)
	limit, _ := strconv.ParseUint(limitStr, 10, 64)

	filters := concertsModel.Filters{Year: year, ComedianId: comedianId, SortBy: sortBy, ExcludedId: excludedId, Limit: limit}

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
	concertGroup := server.Group("concerts")
	{
		concertGroup.GET("", getConcerts)
	}
	concertGroupAuth := concertGroup.Group("", middleware.AuthMiddleware)
	{
		concertGroupAuth.POST("", middleware.RoleMiddleware([]string{createConcertsPermission}), createConcert)
	}
	currentConcertGroup := concertGroup.Group(":concertId", handleConcertId)
	{
		currentConcertGroup.GET("", getConcert)
		currentConcertGroup.PUT("views/inc", incViews)
	}
	currentConcertGroupAuth := currentConcertGroup.Group("", middleware.AuthMiddleware, handleConcertId)
	{
		currentConcertGroupAuth.DELETE("", deleteConcert)
		currentConcertGroupAuth.PUT("", updateConcert)
		currentConcertGroupAuth.PUT("like", like)
	}

	server.Static(os.Getenv("PUBLIC_CONCERTS_VIDEOS_DIR"), os.Getenv("CONCERTS_VIDEOS_DIR"))
	server.Static(os.Getenv("PUBLIC_CONCERTS_FRAMES_DIR"), os.Getenv("CONCERTS_FRAMES_DIR"))

	controller.InitRoutes(currentConcertGroupAuth)

	return currentConcertGroup
}
