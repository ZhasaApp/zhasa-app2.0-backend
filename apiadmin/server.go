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
	getUsersByRoleFunc       repository.GetUsersByRoleFunc
	getUsersWithBranchBrands repository.GetUsersWithBranchBrands
	getUserByIdFunc          repository.GetUserByIdFunc
	getUserBranchFunc        repository.GetUserBranchFunc
	updateUserBrands         repository.UpdateUserBrandsFunc
	updateUserFunc           repository.UpdateUserFunc
	updateUserBranchFunc     repository.UpdateUserBranchFunc
	updateUserRole           repository.UpdateUserRoleFunc

	addDisabledUserFunc repository.AddDisabledUserFunc

	getAllBranchesFunc               repository2.GetAllBranches
	getAllBrandsFunc                 brand.GetAllBrandsFunc
	getUserBrandsFunc                brand.GetUserBrandsFunc
	getFilteredUsersWithBranchBrands repository.GetFilteredUsersWithBranchBrands
	addUserRole                      repository.AddUserRoleFunc
	addUserBranch                    repository.AddUserBranchFunc
}

func NewServer(
	getUserByPhoneFunc repository.GetUserByPhoneFunc,
	createUserFunc repository.CreateUserFunc,
	makeManagerAsUserFunc repository.MakeUserAsManagerFunc,
	getUsersWithoutRolesFunc repository.GetUsersWithoutRolesFunc,
	getUsersByRoleFunc repository.GetUsersByRoleFunc,
	getUsersWithBranchBrands repository.GetUsersWithBranchBrands,
	getUserByIdFunc repository.GetUserByIdFunc,
	getUserBranchFunc repository.GetUserBranchFunc,
	updateUserBrands repository.UpdateUserBrandsFunc,
	updateUserFunc repository.UpdateUserFunc,
	updateUserBranchFunc repository.UpdateUserBranchFunc,
	branchesFunc repository2.GetAllBranches,
	brandsFunc brand.GetAllBrandsFunc,
	addDisabledUserFunc repository.AddDisabledUserFunc,
	getUserBrandsFunc brand.GetUserBrandsFunc,
	getFilteredUsersWithBranchBrands repository.GetFilteredUsersWithBranchBrands,
	addUserRole repository.AddUserRoleFunc,
	addUserBranch repository.AddUserBranchFunc,
	updateUserRole repository.UpdateUserRoleFunc,
) *Server {
	return &Server{
		getUserByPhoneFunc:               getUserByPhoneFunc,
		createUserFunc:                   createUserFunc,
		makeUserAsManagerFunc:            makeManagerAsUserFunc,
		getUsersWithoutRolesFunc:         getUsersWithoutRolesFunc,
		getUsersByRoleFunc:               getUsersByRoleFunc,
		getUsersWithBranchBrands:         getUsersWithBranchBrands,
		getUserByIdFunc:                  getUserByIdFunc,
		getUserBranchFunc:                getUserBranchFunc,
		updateUserBrands:                 updateUserBrands,
		updateUserFunc:                   updateUserFunc,
		updateUserBranchFunc:             updateUserBranchFunc,
		addDisabledUserFunc:              addDisabledUserFunc,
		getAllBranchesFunc:               branchesFunc,
		getAllBrandsFunc:                 brandsFunc,
		getUserBrandsFunc:                getUserBrandsFunc,
		getFilteredUsersWithBranchBrands: getFilteredUsersWithBranchBrands,
		addUserRole:                      addUserRole,
		addUserBranch:                    addUserBranch,
		updateUserRole:                   updateUserRole,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
