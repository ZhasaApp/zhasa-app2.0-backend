package apiadmin

import (
	"github.com/gin-gonic/gin"
	repository2 "zhasa2.0/branch/repository"
	"zhasa2.0/brand"
	"zhasa2.0/user/repository"
)

type Server struct {
	getUserByPhoneFunc         repository.GetUserByPhoneFunc
	getUsersWithoutRolesFunc   repository.GetUsersWithoutRolesFunc
	createUserFunc             repository.CreateUserFunc
	makeUserAsManagerFunc      repository.MakeUserAsManagerFunc
	getUsersByRoleFunc         repository.GetUsersByRoleFunc
	getUserByIdFunc            repository.GetUserByIdFunc
	getUserBranchFunc          repository.GetUserBranchFunc
	updateUserBranchBrandsFunc repository.UpdateUserBranchBrandsFunc
	updateUserFunc             repository.UpdateUserFunc

	getAllBranchesFunc repository2.GetAllBranches
	getAllBrandsFunc   brand.GetAllBrandsFunc
	getUserBrandsFunc  brand.GetUserBrandsFunc
}

func NewServer(
	getUserByPhoneFunc repository.GetUserByPhoneFunc,
	createUserFunc repository.CreateUserFunc,
	makeManagerAsUserFunc repository.MakeUserAsManagerFunc,
	getUsersWithoutRolesFunc repository.GetUsersWithoutRolesFunc,
	getUsersByRoleFunc repository.GetUsersByRoleFunc,
	getUserByIdFunc repository.GetUserByIdFunc,
	getUserBranchFunc repository.GetUserBranchFunc,
	updateUserBranchBrandsFunc repository.UpdateUserBranchBrandsFunc,
	updateUserFunc repository.UpdateUserFunc,
	branchesFunc repository2.GetAllBranches,
	brandsFunc brand.GetAllBrandsFunc,
	getUserBrandsFunc brand.GetUserBrandsFunc,
) *Server {
	return &Server{
		getUserByPhoneFunc:         getUserByPhoneFunc,
		createUserFunc:             createUserFunc,
		makeUserAsManagerFunc:      makeManagerAsUserFunc,
		getUsersWithoutRolesFunc:   getUsersWithoutRolesFunc,
		getUsersByRoleFunc:         getUsersByRoleFunc,
		getUserByIdFunc:            getUserByIdFunc,
		getUserBranchFunc:          getUserBranchFunc,
		updateUserBranchBrandsFunc: updateUserBranchBrandsFunc,
		updateUserFunc:             updateUserFunc,
		getAllBranchesFunc:         branchesFunc,
		getAllBrandsFunc:           brandsFunc,
		getUserBrandsFunc:          getUserBrandsFunc,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
