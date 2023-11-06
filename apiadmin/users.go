package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"zhasa2.0/base"
	"zhasa2.0/brand"
	"zhasa2.0/user/entities"
)

type GetAllUsersRequest struct {
	Page     int32  `json:"page" form:"page"`
	PageSize int32  `json:"size" form:"size"`
	RoleKey  string `json:"role_key" form:"role_key"`
}

type GetAllUsersResponse struct {
	Result  []entities.UserWithBrands `json:"result"`
	HasNext bool                      `json:"has_next"`
	Count   int32                     `json:"count"`
}

func (s *Server) GetAllUsers(ctx *gin.Context) {
	var req GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, total, err := s.getUsersByRoleFunc(req.RoleKey, base.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}

	hasNext := total > req.PageSize*(req.Page+1)

	ctx.JSON(http.StatusOK, GetAllUsersResponse{
		Result:  users,
		HasNext: hasNext,
		Count:   total,
	})
}

type GetUsersWithoutRolesRequest struct {
	Search string `json:"search" form:"search" required:"false"`
}

type GetUsersWithoutRolesResponse struct {
	Result []entities.BaseUser `json:"result"`
}

func (s *Server) GetUsersWithoutRoles(ctx *gin.Context) {
	var req GetUsersWithoutRolesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := s.getUsersWithoutRolesFunc(req.Search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, GetUsersWithoutRolesResponse{
		Result: users,
	})
}

func (s *Server) GetAllUsersForm(ctx *gin.Context) {
	users, _, err := s.getUsersByRoleFunc("sales_manager", base.Pagination{
		Page:     0,
		PageSize: 500,
	})
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "users-list.html", gin.H{
		"users": users,
	})
}

func (s *Server) EditUserForm(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 32)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	user, err := s.getUserByIdFunc(int32(userId))
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	userBranch, err := s.getUserBranchFunc(int32(userId))
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	userBrands, err := s.getUserBrandsFunc(int32(userId))
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	branches, err := s.getAllBranchesFunc()
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	brands, err := s.getAllBrandsFunc()
	if err != nil {
		ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	allBrands := make([]brand.SelectedBrand, 0, len(brands))

	for _, userBrand := range brands {
		selected := false
		if userBrands != nil {
			for _, userBrandItem := range userBrands {
				if userBrandItem.Id == userBrand.Id {
					selected = true
				}
			}
		}
		allBrands = append(allBrands, brand.SelectedBrand{
			Brand:    userBrand,
			Selected: selected,
		})
	}

	ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
		"user":       user,
		"userBranch": userBranch,
		"userBrands": userBrands,
		"branches":   branches,
		"brands":     allBrands,
	})
}

func (s *Server) PerformEditUserFromForm(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	branchIdQuery := ctx.PostForm("branch")
	brandIdsQuery := ctx.PostFormArray("brand")

	firstNameForm := ctx.PostForm("first_name")
	lastNameForm := ctx.PostForm("last_name")
	phone := ctx.PostForm("phone")
	phone = strings.ReplaceAll(phone, " ", "")

	userId, err := strconv.ParseInt(userIdParam, 10, 32)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

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
		ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
			"errors": errors,
		})
	}

	user, err := s.getUserByIdFunc(int32(userId))
	if err != nil {
		ctx.HTML(http.StatusOK, "create-user-form.html", gin.H{
			"errors": []string{"user does not exists"},
		})
		return
	}

	if user.FirstName != firstName.String() ||
		user.LastName != lastName.String() ||
		user.Phone.String() != validatedPhone.String() {
		err = s.updateUserFunc(int32(userId), *firstName, *lastName, *validatedPhone)
		if err != nil {
			ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
				"errors": []string{err.Error()},
			})
			return
		}
	}

	branchId, err := strconv.Atoi(branchIdQuery)
	if err != nil {
		ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	brandsIds := make([]int32, len(brandIdsQuery))
	for i, str := range brandIdsQuery {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
				"errors": []string{err.Error()},
			})
			return
		}
		brandsIds[i] = int32(num)
	}

	branch, err := s.getUserBranchFunc(int32(userId))
	if err != nil {
		ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	if branch.ID != int32(branchId) {
		err = s.updateUserBranchFunc(int32(userId), int32(branchId))
		if err != nil {
			ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
				"errors": []string{err.Error()},
			})
			return
		}
	}

	err = s.updateUserBrands(int32(userId), brandsIds)
	if err != nil {
		ctx.HTML(http.StatusOK, "edit-user.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	ctx.HTML(http.StatusOK, "success-page-edit.html", gin.H{
		"userId": userId,
	})
}
