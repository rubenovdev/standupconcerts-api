package controller

import (
	"comedians/src/common"
	"comedians/src/core/concerts/comments/controller"
	"comedians/src/core/concerts/service"
	"comedians/src/core/usersConcerts/model"
	"comedians/src/middleware"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var concertsPath string

const (
	concertsDir              = "concerts"
	createConcertsPermission = "create_concerts"
)

func init() {
	concertsPath = os.Getenv("ROOT_DIR") + "/" + concertsDir
}

func createConcert(c *gin.Context) {
	c.MultipartForm()
	// file, handler, err := c.Request.FormFile("concert")

	form, err := c.MultipartForm()
	log.Print(form)

	if err != nil {
		log.Panic(err)

		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect file"})
	}

	file := form.File["concert"][0]
	openedFile, _ := file.Open()

	filepath, err := service.UploadConcertFile(openedFile, file.Filename)

	if err != nil {
		log.Panic(err)
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{Message: "error uploading file"})

		return
	}

	concert := model.Concert{}
	concert.Filepath = filepath
	concert.Title = c.Request.Form.Get("title")
	userId, _ := c.Get("userId")

	concert.UserId = userId.(uint64)

	if err := service.CreateConcert(concert); err != nil {
		log.Panic(err)

		service.DeleteConcertFile(filepath)
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})

		return
	}

	common.NewResultResponse(c, http.StatusCreated, common.ResultResponse{Result: gin.H{
		"filepath": filepath,
	}})
}

func getConcerts(c *gin.Context) {
	concerts, err := service.GetConrerts()

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
		return
	}

	common.NewResultResponse(c, http.StatusCreated, common.ResultResponse{Result: concerts})
}

func deleteConcert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
	}

	err = service.DeleteConcert(id)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getConcert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
	}

	concert, err := service.GetConcert(id)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: concert})
}

func updateConcert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
	}

	var concert model.Concert

	if err := c.ShouldBindJSON(&concert); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert data"})
	}

	concert.Id = id

	err = service.UpdateConcert(concert)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func InitGroup(server *gin.Engine) *gin.RouterGroup {
	group := server.Group("concerts", middleware.AuthMiddleware)

	concertPublicPath := os.Getenv("PUBLIC_ROOT_DIR") + "/" + concertsDir

	server.Static(concertPublicPath, concertsPath)

	group.POST("/", middleware.RoleMiddleware([]string{createConcertsPermission}), createConcert)
	group.GET("/", getConcerts)

	group.DELETE("/:concertId", deleteConcert)
	group.PUT("/:concertId", updateConcert)
	group.GET("/:concertId", getConcert)

	controller.InitRoutes(group)

	return group
}
