package apiadmin

import (
	"github.com/gin-gonic/gin"
	repository2 "zhasa2.0/branch/repository"
	"zhasa2.0/brand"
	"zhasa2.0/user/repository"
)

type Server struct {
	getUserByPhoneFunc       repository.GetUserByPhoneFunc
	getUsersWithoutRolesFunc repository.GetUsersWithoutRolesFunc
	createUserFunc           repository.CreateUserFunc
	makeUserAsManagerFunc    repository.MakeUserAsManagerFunc

	getAllBranchesFunc repository2.GetAllBranches
	getAllBrandsFunc   brand.GetAllBrandsFunc
}

func NewServer(
	getUserByPhoneFunc repository.GetUserByPhoneFunc,
	createUserFunc repository.CreateUserFunc,
	makeManagerAsUserFunc repository.MakeUserAsManagerFunc,
	getUsersWithoutRolesFunc repository.GetUsersWithoutRolesFunc,
	branchesFunc repository2.GetAllBranches,
	brandsFunc brand.GetAllBrandsFunc,
) *Server {
	return &Server{
		getUserByPhoneFunc:       getUserByPhoneFunc,
		createUserFunc:           createUserFunc,
		makeUserAsManagerFunc:    makeManagerAsUserFunc,
		getUsersWithoutRolesFunc: getUsersWithoutRolesFunc,
		getAllBranchesFunc:       branchesFunc,
		getAllBrandsFunc:         brandsFunc,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
