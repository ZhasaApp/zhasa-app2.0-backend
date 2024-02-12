package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"zhasa2.0/base"
	"zhasa2.0/brand"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/user/entities"
	"zhasa2.0/user/service"
)

type GetAllUsersByRoleRequest struct {
	Page     int32  `json:"page" form:"page"`
	PageSize int32  `json:"size" form:"size"`
	RoleKey  string `json:"role_key" form:"role_key"`
	Search   string `json:"search" form:"search"`
}

type GetAllUsersByRoleResponse struct {
	Result  []entities.UserWithBrands `json:"result"`
	HasNext bool                      `json:"has_next"`
	Count   int32                     `json:"count"`
}

func (s *Server) GetAllUsersByRole(ctx *gin.Context) {
	var req GetAllUsersByRoleRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pagination := base.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	users, total, err := s.getUsersByRoleFunc(req.Search, req.RoleKey, pagination)
	if err != nil {
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}

	hasNext := pagination.HasNext(total)

	ctx.JSON(http.StatusOK, GetAllUsersByRoleResponse{
		Result:  users,
		HasNext: hasNext,
		Count:   total,
	})
}

type GetAllUsersRequest struct {
	Page      int32    `json:"page" form:"page"`
	PageSize  int32    `json:"size" form:"size"`
	Roles     []string `json:"role_keys" form:"role_keys"`
	Brands    []int32  `json:"brand_ids" form:"brand_ids"`
	Branches  []int32  `json:"branch_ids" form:"branch_ids"`
	Search    string   `json:"search" form:"search"`
	SortType  string   `json:"sort_type" form:"sort_type"`
	SortField string   `json:"sort_field" form:"sort_field"`
}

type GetAllUsersResponse struct {
	Result  []entities.UserWithBranchBrands `json:"result"`
	HasNext bool                            `json:"has_next"`
	Count   int32                           `json:"count"`
}

func (s *Server) GetAllUsers(ctx *gin.Context) {
	var req GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pagination := base.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	params := generated.GetFilteredUsersWithBranchRolesBrandsParams{
		Search:    req.Search,
		Limit:     pagination.PageSize,
		Offset:    pagination.GetOffset(),
		RoleKeys:  req.Roles,
		BrandIds:  req.Brands,
		BranchIds: req.Branches,
		SortField: req.SortField,
		SortType:  req.SortType,
	}
	users, total, err := s.getFilteredUsersWithBranchBrands(params)
	if err != nil {
		ctx.JSON(http.StatusOK, errorResponse(err))
		return
	}

	hasNext := pagination.HasNext(total)

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
	users, _, err := s.getUsersByRoleFunc("", "sales_manager", base.Pagination{
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
		brId := int32(branchId)
		err = s.updateUserBranchFunc(int32(userId), &brId)
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

func (s *Server) DisableUserForm(ctx *gin.Context) {
	userIdParam := ctx.Param("id")
	userId, err := strconv.ParseInt(userIdParam, 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = s.getUserByIdFunc(int32(userId))
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid user id: %d", userId)
		return
	}

	err = s.addDisabledUserFunc(int32(userId))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Redirect(http.StatusSeeOther,
		"/users/all")
}

type AdminLoginRequest struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

func (s *Server) AdminLogin(ctx *gin.Context) {
	var req AdminLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	phone, err := entities.NewPhone(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := s.authService.AdminLogin(*phone, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := s.tokenService.GenerateToken(&service.UserTokenData{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone.String(),
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, AdminLoginResponse{
		Token: string(token),
	})
}
