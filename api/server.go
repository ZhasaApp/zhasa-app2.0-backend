package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	"zhasa2.0/apiadmin"
	. "zhasa2.0/branch/repository"
	. "zhasa2.0/branch_director/repo"
	. "zhasa2.0/brand"
	. "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	. "zhasa2.0/news/repository"
	. "zhasa2.0/owner/repository"
	"zhasa2.0/rating"
	. "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic/repository"
	"zhasa2.0/user/entities"
	. "zhasa2.0/user/repository"
	"zhasa2.0/user/service"
)

type Server struct {
	router *gin.Engine
	apiadmin.Server
	tokenService                                  service.TokenService
	authService                                   service.AuthorizationService
	rankRepo                                      RankingsRepository
	postRepository                                PostRepository
	ownerRepository                               OwnerRepository
	saleTypeRepo                                  SaleTypeRepository
	directorRepo                                  BranchDirectorRepository
	saleRepo                                      SaleRepository
	userBrandGoal                                 UserBrandGoalFunc
	getUserBrandFunc                              GetUserBrandFunc
	updateUserBrandRatio                          UpdateUserBrandRatioFunc
	getUserRatingFunc                             rating.GetUserRatingFunc
	userRepo                                      UserRepository
	getUserBranchFunc                             GetUserBranchFunc
	calculateUserBrandRatio                       CalculateUserBrandRatio
	getBranchBrands                               GetBranchBrandsFunc
	getAllBrands                                  GetAllBrandsFunc
	getUserBrands                                 GetUserBrandsFunc
	getBranchBrandFunc                            GetBranchBrandFunc
	getBranchBrandSaleSumFunc                     GetBranchBrandSaleSumFunc
	getBranchBrandGoalFunc                        GetBranchBrandGoalFunc
	getUsersOrderedByRatioForGivenBrandFunc       GetUsersOrderedByRatioForGivenBrandFunc
	getBranchUsersOrderedByRatioForGivenBrandFunc GetBranchUsersOrderedByRatioForGivenBrandFunc
	getBranchByIdFunc                             GetBranchByIdFunc
	getBranchesByBrandFunc                        GetBranchesByBrandFunc
	setBranchBrandSaleTypeGoal                    SetBranchBrandSaleTypeGoal
	setUserBrandGoalRequest                       SetUserBrandSaleTypeGoalFunc
	getUserByBranchBrandRoleFunc                  GetUserByBranchBrandRoleFunc
	getBranchBrandMonthlyYearStatisticFunc        GetBranchBrandMonthlyYearStatisticFunc
	getUsersByBranchBrandRoleFunc                 GetUsersByBranchBrandRoleFunc
	getSaleSumByUserBrandTypePeriodFunc           GetSaleSumByUserBrandTypePeriodFunc
	salesByBrandUserFunc                          SalesByBrandUserFunc
	saleAddFunc                                   SaleAddFunc
	saleEditFunc                                  SaleEditFunc
	ratedBranchesFunc                             RatedBranchesFunc

	// user functions
	createUserFunc     CreateUserFunc
	getUserByPhoneFunc GetUserByPhoneFunc
	getUserByIdFunc    GetUserByIdFunc
	uploadAvatarFunc   UploadAvatarFunc
	deleteAvatarFunc   DeleteAvatarFunc
}

func (server *Server) InitSuperUser() error {
	request := entities.CreateUserRequest{
		FirstName: "admin",
		LastName:  "admin",
		Phone:     "+77081070480",
	}
	_, err := server.getUserByPhoneFunc(request.Phone)

	if err == nil {
		fmt.Println("super user already exist")
		return nil
	}

	_, err = server.createUserFunc(request.FirstName, request.LastName, request.Phone)
	if err != nil {
		return err
	}
	fmt.Println("super user created")
	return nil
}

func NewServer(ctx context.Context) *Server {
	server := &Server{}
	initDependencies(server, ctx)

	router := gin.Default()

	router.POST("/image/avatar/upload", verifyToken(server.tokenService), server.HandleAvatarUpload)
	router.POST("/image/news/upload", verifyToken(server.tokenService), server.HandleNewsUpload)

	router.POST("/user/avatar", verifyToken(server.tokenService), server.UploadUserAvatar)
	router.DELETE("/user/avatar", verifyToken(server.tokenService), server.DeleteAvatar)
	router.POST("/csv/managers", verifyToken(server.tokenService), server.HandleManagersUpload)
	router.POST("/csv/directors", verifyToken(server.tokenService), server.HandleDirectorsUpload)

	authRoute := router.Group("auth/")
	{
		authRoute.POST("/request-code", server.requestAuthCode)
		authRoute.POST("/try", server.tryAuth)
	}

	router.GET("user/get-user-profile", server.getUserProfile)

	adminRoute := router.Group("admin/").Use(verifyToken(server.tokenService))
	{
		adminRoute.POST("/user", server.CreateUser)
		adminRoute.POST("/manager", server.CreateManager)
		adminRoute.GET("/branches", server.GetAllBranches)
		adminRoute.GET("/brands", server.GetAllBrands)
	}

	smRoute := router.Group("sales-manager/")
	smRoute.GET("/year-statistic", server.GetUserBrandYearStatistic).Use(verifyToken(server.tokenService))
	smRoute.GET("/sale/list", server.GetSales).Use(verifyToken(server.tokenService))

	branchRoute := router.Group("branch/").Use(verifyToken(server.tokenService))
	{
		branchRoute.GET("/year-statistic", server.GetBranchBrandYearStatistic)
		branchRoute.GET("/sales-managers", server.GetBranchSalesManagerList)
	}

	router.DELETE("sales/delete", server.DeleteSale).Use(verifyToken(server.tokenService))
	router.POST("sales-manager/sale/new", server.AddSale).Use(verifyToken(server.tokenService))
	router.POST("sales/edit", server.EditSale).Use(verifyToken(server.tokenService))

	directorRouter := router.Group("director/")
	{
		directorRouter.POST("sales-manager/goal", server.SetUserBrandGoal).Use(verifyToken(server.tokenService))
		directorRouter.GET("sales-manager/goal", server.GetSmGoal).Use(verifyToken(server.tokenService))
		directorRouter.POST("branch/goal", server.SetBranchGoal).Use(verifyToken(server.tokenService))
		directorRouter.GET("branch/goal", server.GetBranchGoal).Use(verifyToken(server.tokenService))
	}

	router.GET("sales-manager/dashboard", server.SMDashboard).Use(verifyToken(server.tokenService))

	router.GET("branch/dashboard", server.BranchDashboard).Use(verifyToken(server.tokenService))
	router.GET("rating/branches", server.GetRatedBranches).Use(verifyToken(server.tokenService))
	router.GET("rating/managers", server.GetOrderedUsers).Use(verifyToken(server.tokenService))

	router.GET("news", verifyToken(server.tokenService), server.GetPosts)
	router.POST("news/new", verifyToken(server.tokenService), server.CreatePost)
	router.DELETE("news/delete", verifyToken(server.tokenService), server.DeletePost)
	router.POST("news/like-toggle", verifyToken(server.tokenService), server.ToggleLike)

	router.GET("news/comments", verifyToken(server.tokenService), server.GetComments)
	router.POST("news/comments/new", verifyToken(server.tokenService), server.CreateComment)
	router.DELETE("news/comments/delete", verifyToken(server.tokenService), server.DeleteComment)

	router.GET("user/brands", verifyToken(server.tokenService), server.GetUserBrands)
	router.GET("branch/brands", verifyToken(server.tokenService), server.GetBranchBrands)
	router.GET("brands", verifyToken(server.tokenService), server.GetAllBrands)

	server.router = router
	return server
}

func initDependencies(server *Server, ctx context.Context) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DATA_BASE_URL")
	conn, err := sql.Open(dbDriver, dbSource)

	log.Println(dbDriver, dbSource)
	if err != nil {
		log.Panic("Cannot connect to db", err)
	}

	store := generated.NewStore(conn)
	customQuerier := NewCustomQuerier(conn)
	userRepo := NewUserRepository(ctx, store)
	saleTypeRepo := NewSaleTypeRepository(ctx, store)
	branchRepo := NewBranchRepository(ctx, store, customQuerier, saleTypeRepo)
	directorRepo := NewBranchDirectorRepository(ctx, store)
	rankingsRepo := NewRankingsRepository(ctx, customQuerier, branchRepo)
	postRepo := NewPostRepository(ctx, store, customQuerier)
	ownerRepo := NewOwnerRepo(ctx, store)
	authService := service.NewAuthorizationService(ctx, userRepo)
	brandGoal := NewUserGoalFunc(ctx, store)
	userSaleSum := NewGetSaleSumByUserBrandTypePeriodFunc(ctx, store)
	saleRepo := NewSaleRepo(ctx, store, saleTypeRepo, brandGoal, userSaleSum)
	allBrands := NewGetAllBrandsFunc(ctx, store)
	encKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

	tokenService := service.NewTokenService(&encKey)
	server.tokenService = tokenService
	server.authService = authService
	server.postRepository = postRepo
	server.ownerRepository = ownerRepo
	server.directorRepo = directorRepo
	server.rankRepo = rankingsRepo
	server.saleRepo = saleRepo
	server.saleTypeRepo = saleTypeRepo
	server.userBrandGoal = brandGoal
	server.getUserBrandFunc = NewGetUserBrandFunc(ctx, store)
	server.updateUserBrandRatio = NewUpdateUserBrandRatioFunc(ctx, store)
	server.getUserRatingFunc = rating.NewGetUserRatingFunc(ctx, store)
	server.userRepo = userRepo
	server.getUserBranchFunc = NewGetUserBranchFunc(ctx, store)
	server.calculateUserBrandRatio = NewCalculateUserBrandRatio(saleTypeRepo, userSaleSum, server.userBrandGoal)
	server.getBranchBrands = NewGetBranchBrandsFunc(ctx, store)
	server.getAllBrands = allBrands
	server.getUserBrands = NewGetUserBrandsFunc(ctx, store)
	server.getBranchBrandFunc = NewGetBranchBrand(ctx, store)
	server.getBranchBrandSaleSumFunc = NewGetBranchBrandSaleSumFunc(ctx, store)
	server.getBranchBrandGoalFunc = NewGetBranchBrandGoalFunc(ctx, store)
	server.getUsersOrderedByRatioForGivenBrandFunc = NewGetUsersOrderedByRatioForGivenBrandFunc(ctx, store)
	server.getBranchUsersOrderedByRatioForGivenBrandFunc = NewGetUsersOrderedBYRatioForGivenBrandAndBranchFunc(ctx, store)
	server.getBranchByIdFunc = NewGetBranchByIdFunc(ctx, store)
	server.getBranchesByBrandFunc = NewGetBranchesByBrandFunc(ctx, store)
	server.setBranchBrandSaleTypeGoal = NewSetBranchGoalFunc(ctx, store)
	server.setUserBrandGoalRequest = NewSetUserBrandSaleTypeGoalFunc(ctx, store)
	server.getUserByBranchBrandRoleFunc = NewGetUserByBranchBrandRoleFunc(ctx, store)
	server.getBranchBrandMonthlyYearStatisticFunc = NewGetBranchBrandMonthlyYearStatisticFunc(saleTypeRepo, server.getBranchBrandGoalFunc, server.getBranchBrandFunc, server.getBranchBrandSaleSumFunc)
	server.getUsersByBranchBrandRoleFunc = NewGetUsersByBranchBrandRoleFunc(ctx, store)
	server.getSaleSumByUserBrandTypePeriodFunc = userSaleSum
	server.salesByBrandUserFunc = NewSalesByBrandUserFunc(ctx, store)
	server.saleAddFunc = NewSaleAddFunc(ctx, store)
	server.saleEditFunc = NewSaleEditFunc(ctx, store)

	server.ratedBranchesFunc = NewRatedBranchesFunc(ctx, store, server.getBranchBrandSaleSumFunc, server.getBranchBrandGoalFunc)

	// user functions
	server.createUserFunc = NewCreateUserFunc(ctx, store)
	server.getUserByPhoneFunc = NewGetUserByPhoneFunc(ctx, store)
	server.getUserByIdFunc = NewGetUserByIdFunc(ctx, store)
	server.uploadAvatarFunc = NewUploadAvatarFunc(ctx, store)
	server.deleteAvatarFunc = NewDeleteAvatarFunc(ctx, store)

	getUserByPhoneFunc := NewGetUserByPhoneFunc(ctx, store)
	createUserFunc := NewCreateUserFunc(ctx, store)
	makeUserAsManagerFunc := NewMakeUserAsManagerFunc(ctx, store)
	getAllBranches := NewGetAllBranchesFunc(ctx, store)

	server.Server = *apiadmin.NewServer(getUserByPhoneFunc, createUserFunc, makeUserAsManagerFunc, getAllBranches, allBrands)
}

// Start runs the HTTP server a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
