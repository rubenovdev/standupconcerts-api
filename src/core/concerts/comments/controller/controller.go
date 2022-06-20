package controller

import (
	"comedians/src/common"
	"comedians/src/core/concerts/comments/model"
	"comedians/src/core/concerts/comments/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createComment(c *gin.Context) {
	userId, _ := c.Get("userId")
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	var input model.CreateCommentDto

	if err := c.BindJSON(&input); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect data"))

		return
	}

	comment := model.Comment{}

	comment.ConcertId = concertId
	comment.Body = input.Body
	comment.UserId = userId.(uint64)

	if err := service.CreateComment(comment); err != nil {
		c.Status(500)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func getComments(c *gin.Context) {
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)
	comments, err := service.GetComments(concertId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{Result: gin.H{
		"comments": comments,
	}})

}

func deleteComment(c *gin.Context) {
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)
	commentId, _ := strconv.ParseUint(c.Param("commentId"), 10, 64)

	err := service.DeleteComment(concertId, commentId)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func updateComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)

	var concert model.Comment

	if err := c.ShouldBindJSON(&concert); err != nil {
		common.NewErrorResponse(c, http.StatusBadRequest, errors.New("incorrect concert data"))
	}

	concert.Id = id

	err := service.UpdateComment(concert)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func like(c *gin.Context) {
	concertId, _ := strconv.ParseUint(c.Param("concertId"), 10, 64)
	commentId, _ := strconv.ParseUint(c.Param("commentId"), 10, 64)
	userId, _ := c.Get("userId")

	if err := service.Like(concertId, commentId, userId.(uint64)); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	common.NewResultResponse(c, http.StatusOK, common.ResultResponse{})
}

func handleCommentId(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("commentId"), 10, 64)

	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, errors.New("incorrect comment id"))

		return
	}
}

func InitRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	commentsGroup := router.Group("comments")
	{
		commentsGroup.POST("", createComment)
		commentsGroup.GET("", getComments)
	}

	currentCommentGroup := commentsGroup.Group(":commentId", handleCommentId)

	{
		currentCommentGroup.DELETE("", deleteComment)
		currentCommentGroup.PUT("", updateComment)
		currentCommentGroup.PUT("like", like)
	}

	return commentsGroup
}
