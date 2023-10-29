package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	. "zhasa2.0/api/entities"
	. "zhasa2.0/sale/entities"
	entities2 "zhasa2.0/user/entities"
)

func (server *Server) getSaleTypes(ctx *gin.Context) {
	types, err := server.saleTypeRepo.GetSaleTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	saleTypes := make([]SaleTypeResponse, 0)
	for _, item := range *types {
		saleTypes = append(saleTypes, SaleTypeResponse{
			Id:    int32(item.Id),
			Title: item.Title,
			Color: item.Color,
		})
	}

	saleTypesResponse := SaleTypesResponse{
		Result: saleTypes,
	}

	ctx.JSON(http.StatusOK, saleTypesResponse)
}

func (server *Server) createSaleType(ctx *gin.Context) {
	var createSaleTypeBody CreateSaleTypeRequest
	if err := ctx.ShouldBindJSON(&createSaleTypeBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	body := CreateSaleTypeBody{
		Title:       createSaleTypeBody.Title,
		Description: createSaleTypeBody.Description,
	}

	id, err := server.saleTypeRepo.CreateSaleType(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, IdResponse{
		Id: int32(id),
	})
}

func (server *Server) createBranchDirector(ctx *gin.Context) {
	var body CreateBranchDirectorBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := entities2.NewName(body.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := entities2.NewName(body.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	phone, err := entities2.NewPhone(body.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.getUserByPhoneFunc(*phone)
	if user != nil && err == nil {

		//id, err := server.directorService.CreateBranchDirector(user.Id, body.BranchId)
		//
		//if err != nil {
		//	fmt.Println(err)
		//	ctx.JSON(http.StatusBadRequest, errors.New("create branch director error"))
		//	return
		//}
		//ctx.JSON(http.StatusOK, IdResponse{
		//	Id: int32(id),
		//})
		return
	}

	_, err = server.createUserFunc(*firstName, *lastName, *phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err = server.getUserByPhoneFunc(*phone)
	if user == nil && err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//	id, err := server.directorService.CreateBranchDirector(user.Id, body.BranchId)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errors.New("create branch director error"))
		return
	}
	//response := IdResponse{
	//	Id: int32(id),
	//}
	//ctx.JSON(http.StatusOK, response)
}
