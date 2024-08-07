package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetGoodsByBrandRequest struct {
	BrandId int32 `json:"brand_id" form:"brand_id" binding:"required"`
}

type GoodItem struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GoodsResponse struct {
	Result []GoodItem `json:"result"`
}

func (server *Server) GetGoodsByBrand(ctx *gin.Context) {
	var getGoodsByBrandRequest GetGoodsByBrandRequest
	if err := ctx.ShouldBindQuery(&getGoodsByBrandRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	goods, err := server.getGoodByBrandFunc(getGoodsByBrandRequest.BrandId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	goodItems := make([]GoodItem, 0)

	for _, good := range goods {
		goodItems = append(goodItems, GoodItem{
			Id:          good.Id,
			Name:        good.Name,
			Description: good.Description,
		})
	}

	response := GoodsResponse{
		Result: goodItems,
	}

	ctx.JSON(http.StatusOK, response)
}
