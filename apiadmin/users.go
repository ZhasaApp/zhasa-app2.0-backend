package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"zhasa2.0/brand"
	"zhasa2.0/user/entities"
)

func (s *Server) GetAllUsersForm(ctx *gin.Context) {
	users, err := s.getUsersByRoleFunc("sales_manager")
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
