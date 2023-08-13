package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/base"
)

type getCommentsRequest struct {
	Page     int32 `json:"page" form:"page"`
	PageSize int32 `json:"limit" form:"limit"`
	NewsId   int32 `json:"news_id" form:"news_id"`
}

type CommentItem struct {
	Id      int32  `json:"id"`
	Author  Author `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type CommentsResponse struct {
	Result  []CommentItem `json:"result"`
	HasNext bool          `json:"has_next"`
	Count   int32         `json:"count"`
}

func (server Server) GetComments(ctx *gin.Context) {
	var req *getCommentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := server.postRepository.GetPostComments(req.NewsId, base.Pagination{
		PageSize: req.PageSize,
		Page:     req.Page,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	comments := make([]CommentItem, 0)

	for _, row := range data {
		comments = append(comments, CommentItem{
			Author: Author{
				Id:       row.User.Id,
				Avatar:   row.User.AvatarPointer(),
				FullName: row.User.GetFullName(),
			},
			Id:      row.Id,
			Message: row.Message,
			Date:    row.CreatedDate.Format("2006-01-02 15:04:05"),
		})
	}

	count := int32(len(comments))

	hasNext := count > req.PageSize*(req.Page+1)

	ctx.JSON(http.StatusOK, CommentsResponse{
		Result:  comments,
		HasNext: hasNext,
		Count:   count,
	})
}
