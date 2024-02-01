package apiadmin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"zhasa2.0/user/entities"
)

type CreateUserRequest struct {
	Phone     string  `json:"phone"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	RoleKey   string  `json:"role"`
	Brands    []int32 `json:"brand_ids"`
	BranchID  *int32  `json:"branch_id"`
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

	if err = s.updateUserBrands(id, request.Brands); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if request.BranchID != nil {
		if err = s.addUserBranch(id, *request.BranchID); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if err = s.addUserRole(id, request.RoleKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (s *Server) GetUserForm(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create-user-form.html", gin.H{})
}

func (s *Server) CreateUserFromForm(ctx *gin.Context) {
	firstNameForm := ctx.PostForm("first_name")
	lastNameForm := ctx.PostForm("last_name")
	phone := ctx.PostForm("phone")
	phone = strings.ReplaceAll(phone, " ", "")

	var errors []string

	firstName, err := entities.NewName(firstNameForm)
	if err != nil {
		errors = append(errors, err.Error())
	}

	lastName, err := entities.NewName(lastNameForm)
	if err != nil {
		errors = append(errors, err.Error())
	}

	validatedPhone, err := entities.NewPhone(phone)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		ctx.HTML(http.StatusOK, "create-user-form.html", gin.H{
			"errors": errors,
		})
	}

	user, err := s.getUserByPhoneFunc(*validatedPhone)
	if user != nil {
		ctx.HTML(http.StatusOK, "create-user-form.html", gin.H{
			"errors": []string{"user already exists for the given phone number"},
		})
		return
	}

	id, err := s.createUserFunc(*firstName, *lastName, *validatedPhone)
	if err != nil {
		ctx.HTML(http.StatusOK, "create-user-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	ctx.Redirect(http.StatusSeeOther,
		"/create-manager?id="+strconv.Itoa(int(id)))
}
