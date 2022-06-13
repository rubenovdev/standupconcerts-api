package controller

import (
	"comedians/src/common"
	"comedians/src/core/concerts/comments/model"
	"comedians/src/core/concerts/comments/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createComment(c *gin.Context) {
	var input model.CreateCommentDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect data"})

		return
	}

	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
		return
	}

	comment := model.Comment{}

	comment.ConcertId = concertId
	comment.Body = input.Body
	comment.UserId = input.UserId

	if err := service.CreateComment(comment); err != nil {
		c.Status(500)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getComments(c *gin.Context) {
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
		return
	}

	comments, err := service.GetComments(concertId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: comments})

}

func deleteComment(c *gin.Context) {
	concertId, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
	}

	commentId, err := strconv.ParseUint(c.Param("commentId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect comment id"})
	}

	err = service.DeleteComment(concertId, commentId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func updateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("concertId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert id"})
	}

	var concert model.Comment

	if err := c.ShouldBindJSON(&concert); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorResponse{Message: "incorrect concert data"})
	}

	concert.Id = id

	err = service.UpdateComment(concert)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorResponse{})
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func InitRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	commentsGroup := router.Group("/:concertId/comments")

	commentsGroup.POST("/", createComment)
	commentsGroup.GET("/", getComments)

	commentsGroup.DELETE("/:commentId", deleteComment)
	commentsGroup.PUT("/:commentsId", updateComment)

	return commentsGroup
}
