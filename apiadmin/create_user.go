package apiadmin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhasa2.0/user/entities"
)

type CreateUserRequest struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var request CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	firstName, err := entities.NewName(request.FirstName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lastName, err := entities.NewName(request.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	validatedPhone, err := entities.NewPhone(request.Phone)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.getUserByPhoneFunc(*validatedPhone)
	if user != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("user already exists for the given phone number")))
		return
	}

	id, err := s.createUserFunc(*firstName, *lastName, *validatedPhone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
