package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	entities3 "zhasa2.0/branch/entities"
	"zhasa2.0/sale/entities"
	entities2 "zhasa2.0/user/entities"
)

type createSaleTypeBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type idResponse struct {
	Id int32 `json:"id"`
}

type createBranchDirectorBody struct {
	CreateUserBody
	BranchId int32 `json:"branch_id"`
}

func (server *Server) createSaleType(ctx *gin.Context) {
	var createSaleTypeBody createSaleTypeBody
	if err := ctx.ShouldBindJSON(&createSaleTypeBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	body := entities.CreateSaleTypeBody{
		Title:       createSaleTypeBody.Title,
		Description: createSaleTypeBody.Description,
	}

	id, err := server.saleTypeService.CreateSaleType(body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, idResponse{
		Id: int32(id),
	})
}

func (server *Server) createBranchDirector(ctx *gin.Context) {
	var body createBranchDirectorBody
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

	createUserRequest := entities2.CreateUserRequest{
		Phone:     *phone,
		FirstName: *firstName,
		LastName:  *lastName,
	}

	user, err := server.userService.GetUserByPhone(*phone)
	if user != nil && err == nil {

		id, err := server.directorService.CreateBranchDirector(entities2.UserId(user.Id), entities3.BranchId(body.BranchId))

		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, errors.New("create branch director error"))
			return
		}
		ctx.JSON(http.StatusOK, idResponse{
			Id: int32(id),
		})
		return
	}

	err = server.userService.CreateUser(createUserRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err = server.userService.GetUserByPhone(*phone)
	if user == nil && err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := server.directorService.CreateBranchDirector(entities2.UserId(user.Id), entities3.BranchId(body.BranchId))

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errors.New("create branch director error"))
		return
	}
	response := idResponse{
		Id: int32(id),
	}
	ctx.JSON(http.StatusOK, response)
}
