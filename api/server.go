package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
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
	router                  *gin.Engine
	userService             service.UserService
	tokenService            service.TokenService
	authService             service.AuthorizationService
	rankRepo                RankingsRepository
	postRepository          PostRepository
	ownerRepository         OwnerRepository
	saleTypeRepo            SaleTypeRepository
	directorRepo            BranchDirectorRepository
	saleRepo                SaleRepository
	userBrandGoal           UserBrandGoalFunc
	getUserBrandFunc        GetUserBrandFunc
	updateUserBrandRatio    UpdateUserBrandRatioFunc
	getUserRatingFunc       rating.GetUserRatingFunc
	userRepo                UserRepository
	getUserBranchFunc       GetUserBranchFunc
	calculateUserBrandRatio CalculateUserBrandRatio
	getBranchBrands         GetBranchBrandsFunc
	getAllBrands            GetAllBrandsFunc
	getUserBrands           GetUserBrandsFunc
}

func (server *Server) InitSuperUser() error {
	request := entities.CreateUserRequest{
		FirstName: "admin",
		LastName:  "admin",
		Phone:     "+77081070480",
	}
	_, err := server.userRepo.GetUserByPhone("+77081070480")

	if err == nil {
		fmt.Println("super user already exist")
		return nil
	}

	_, err = server.userService.CreateUser(request)
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
		adminRoute.POST("/user/new", server.createUser)
		adminRoute.POST("/branch/new", server.createBranch)
		adminRoute.POST("/sale-type/new", server.createSaleType)
		adminRoute.POST("/branch-director/new", server.createBranchDirector)
		adminRoute.GET("/branch/list", server.GetBranchList)
		adminRoute.GET("sale-type/list", server.getSaleTypes)
	}

	smRoute := router.Group("sales-manager/")
	smRoute.GET("/branch/list", server.GetBranchList).Use(verifyToken(server.tokenService))
	smRoute.GET("/year-statistic", server.GetUserBrandYearStatistic).Use(verifyToken(server.tokenService))
	smRoute.GET("/sale/list", server.GetSales).Use(verifyToken(server.tokenService))

	branchRoute := router.Group("branch/").Use(verifyToken(server.tokenService))
	{
		branchRoute.GET("/year-statistic", server.getBranchYearStatistic)
		branchRoute.GET("/sales-managers", server.GetBranchSalesManagerList)
	}

	router.DELETE("sales/delete", server.DeleteSale).Use(verifyToken(server.tokenService))
	router.POST("sales/edit", server.EditSale).Use(verifyToken(server.tokenService))
	router.POST("sales/new", server.EditSale).Use(verifyToken(server.tokenService))

	directorRouter := router.Group("director/")
	{
		directorRouter.POST("sales-manager/goal", server.SetSmGoal).Use(verifyToken(server.tokenService))
		directorRouter.GET("sales-manager/goal", server.GetSmGoal).Use(verifyToken(server.tokenService))
		directorRouter.POST("branch/goal", server.SetBranchGoal).Use(verifyToken(server.tokenService))
		directorRouter.GET("branch/goal", server.GetBranchGoal).Use(verifyToken(server.tokenService))
	}

	router.GET("sales-manager/dashboard", server.SMDashboard).Use(verifyToken(server.tokenService))

	router.GET("branch/dashboard", server.getBranchDashboardStatistic).Use(verifyToken(server.tokenService))
	router.GET("rating/branches", server.GetBranchList)
	//	router.GET("rating/managers", server.GetSalesManagers).Use(verifyToken(server.tokenService))

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
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthorizationService(ctx, userRepo)
	brandGoal := NewUserGoalFunc(ctx, store)
	saleRepo := NewSaleRepo(ctx, store, saleTypeRepo, brandGoal)
	encKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

	tokenService := service.NewTokenService(&encKey)
	server.userService = userService
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
	server.calculateUserBrandRatio = NewCalculateUserBrandRatio(saleTypeRepo, saleRepo, server.userBrandGoal, server.getUserBrandFunc)
	server.getBranchBrands = NewGetBranchBrandsFunc(ctx, store)
	server.getAllBrands = NewGetAllBrandsFunc(ctx, store)
	server.getUserBrands = NewGetUserBrandsFunc(ctx, store)
}

// Start runs the HTTP server a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
