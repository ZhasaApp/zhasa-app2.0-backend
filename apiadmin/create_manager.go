package apiadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CreateManagerRequest struct {
	UserId   int32   `json:"user_id"`
	BranchId int32   `json:"branch_id"`
	Brands   []int32 `json:"brands"`
}

func (s *Server) CreateManager(ctx *gin.Context) {
	var request CreateManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.makeUserAsManagerFunc(request.UserId, request.BranchId, request.Brands)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) CreateManagerForm(ctx *gin.Context) {
	var (
		err          error
		selectedUser int
	)
	selectedUserQuery := ctx.Query("id")
	if selectedUserQuery != "" {
		selectedUser, err = strconv.Atoi(selectedUserQuery)
		if err != nil {
			ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
				"errors": []string{err.Error()},
			})
			return
		}
	}

	users, err := s.getUsersWithoutRolesFunc("")
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	branches, err := s.getAllBranchesFunc()
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	brands, err := s.getAllBrandsFunc()
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
		"selectedUser": selectedUser,
		"users":        users,
		"branches":     branches,
		"brands":       brands,
	})
}

func (s *Server) CreateManagerFromForm(ctx *gin.Context) {
	userIdQuery := ctx.PostForm("user")
	branchIdQuery := ctx.PostForm("branch")
	brandIdsQuery := ctx.PostFormArray("brand")

	userId, err := strconv.Atoi(userIdQuery)
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	branchId, err := strconv.Atoi(branchIdQuery)
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	brandsIds := make([]int32, len(brandIdsQuery))
	for i, str := range brandIdsQuery {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
				"errors": []string{err.Error()},
			})
			return
		}
		brandsIds[i] = int32(num)
	}

	err = s.makeUserAsManagerFunc(int32(userId), int32(branchId), brandsIds)
	if err != nil {
		ctx.HTML(http.StatusOK, "create-manager-form.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	ctx.HTML(http.StatusOK, "success-page.html", gin.H{})
}
