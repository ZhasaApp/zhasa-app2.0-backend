package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/base"
)

type NewsListItem struct {
	ID       int32    `json:"id"`
	Author   Author   `json:"author"`
	Images   []string `json:"images"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Date     string   `json:"date"`
	IsLiked  bool     `json:"is_liked"`
	Likes    int32    `json:"likes"`
	Comments int32    `json:"comments"`
}

type NewsListItemResponse struct {
	Result  []NewsListItem `json:"result"`
	HasNext bool           `json:"has_next"`
	Count   int32          `json:"count"`
}

type Author struct {
	Avatar   *string `json:"avatar"`
	FullName string  `json:"full_name"`
	Id       int32   `json:"id"`
}

type GetPostsRequest struct {
	Page     int32 `json:"page" form:"page"`
	PageSize int32 `json:"limit" form:"limit"`
}

func (server Server) GetPosts(ctx *gin.Context) {
	var req *GetPostsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId := int32(ctx.GetInt("user_id"))
	posts, err := server.postRepository.GetPosts(userId, Pagination{
		PageSize: req.PageSize,
		Page:     req.Page,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	news := make([]NewsListItem, 0)
	for _, post := range posts {
		news = append(news, NewsListItem{
			ID: post.Id,
			Author: Author{
				Id:       post.Author.Id,
				Avatar:   post.Author.AvatarPointer(),
				FullName: post.Author.GetFullName(),
			},
			Images:   post.Images,
			Title:    post.Title,
			Body:     post.Body,
			Date:     post.CreatedDate.Format("2006-01-02 15:04:05"),
			IsLiked:  post.IsLiked,
			Likes:    post.LikesCount,
			Comments: post.CommentsCount,
		})
	}

	count := int32(len(news))

	hasNext := count > req.PageSize*(req.Page+1)

	ctx.JSON(http.StatusOK, NewsListItemResponse{
		Result:  news,
		HasNext: hasNext,
		Count:   count,
	})
}
