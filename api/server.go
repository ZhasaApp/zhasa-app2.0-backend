package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	repository3 "zhasa2.0/branch/repository"
	service3 "zhasa2.0/branch/service"
	"zhasa2.0/branch_director/repo"
	service5 "zhasa2.0/branch_director/service"
	. "zhasa2.0/db/hand-made"
	generated "zhasa2.0/db/sqlc"
	repository2 "zhasa2.0/manager/repository"
	service2 "zhasa2.0/manager/service"
	. "zhasa2.0/news/repository"
	"zhasa2.0/sale/repository"
	service4 "zhasa2.0/sale/service"
	. "zhasa2.0/statistic/repository"
	service6 "zhasa2.0/statistic/service"
	"zhasa2.0/user/entities"
	repository4 "zhasa2.0/user/repository"
	"zhasa2.0/user/service"
)

type Server struct {
	router              *gin.Engine
	userService         service.UserService
	salesManagerService service2.SalesManagerService
	branchService       service3.BranchService
	saleTypeService     service4.SaleTypeService
	tokenService        service.TokenService
	authService         service.AuthorizationService
	directorService     service5.BranchDirectorService
	analyticsService    service6.AnalyticsService
	postRepository      PostRepository
}

func (server Server) InitSuperUser() error {
	request := entities.CreateUserRequest{
		FirstName: "admin",
		LastName:  "admin",
		Phone:     "+77081070480",
	}
	_, err := server.userService.GetUserByPhone("+77081070480")

	if err == nil {
		fmt.Println("super user already exist")
		return nil
	}

	err = server.userService.CreateUser(request)
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

	router.POST("/user/avatar", verifyToken(server.tokenService), server.UploadUserAvatar)
	router.DELETE("/user/avatar", verifyToken(server.tokenService), server.DeleteAvatar)

	authRoute := router.Group("auth/")
	{
		authRoute.POST("/request-code", server.requestAuthCode)
		authRoute.POST("/try", server.tryAuth)
	}

	router.GET("user/get-user-profile", server.getUserProfile)

	adminRoute := router.Group("admin/").Use(verifyToken(server.tokenService))
	{
		adminRoute.POST("/user/new", server.createUser)
		adminRoute.POST("/sales-manager/new", server.createSalesManager)
		adminRoute.POST("/branch/new", server.createBranch)
		adminRoute.POST("/sale-type/new", server.createSaleType)
		adminRoute.POST("/branch-director/new", server.createBranchDirector)
		adminRoute.GET("/branch/list", server.GetBranchList)
		adminRoute.GET("sale-type/list", server.getSaleTypes)
	}

	smRoute := router.Group("sales-manager/")
	smRoute.POST("/sale/new", getSalesManager(server.tokenService, server.salesManagerService), server.saveSale)
	smRoute.GET("/branch/list", server.GetBranchList).Use(verifyToken(server.tokenService))
	smRoute.GET("/year-statistic", server.getYearStatistic).Use(verifyToken(server.tokenService))
	smRoute.GET("/sale/list", server.getSales).Use(verifyToken(server.tokenService))

	branchRoute := router.Group("branch/").Use(verifyToken(server.tokenService))
	{
		branchRoute.GET("/year-statistic", server.getBranchYearStatistic)
		branchRoute.GET("/sales-managers", server.GetBranchSalesManagerList)
	}

	router.DELETE("sales/delete", server.DeleteSale).Use(verifyToken(server.tokenService))
	router.POST("sales/edit", server.EditSale).Use(verifyToken(server.tokenService))

	directorRouter := router.Group("director/")
	{
		directorRouter.POST("sales-manager/goal", server.SetSmGoal).Use(getBranchDirector(server.tokenService, server.directorService))
		directorRouter.GET("sales-manager/goal", server.GetSmGoal).Use(getBranchDirector(server.tokenService, server.directorService))
	}

	router.GET("sales-manager/dashboard", server.getSalesManagerDashboardStatistic).Use(verifyToken(server.tokenService))

	router.GET("branch/dashboard", server.getBranchDashboardStatistic).Use(verifyToken(server.tokenService))
	router.GET("rating/branches", server.GetBranchList)
	router.GET("rating/managers", server.GetSalesManagers).Use(verifyToken(server.tokenService))

	router.GET("news", verifyToken(server.tokenService), server.GetPosts)
	router.GET("news/new", verifyToken(server.tokenService), server.CreatePost)
	router.GET("news/delete", verifyToken(server.tokenService), server.DeletePost)
	server.router = router
	return server
}

func initDependencies(server *Server, ctx context.Context) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbSource := os.Getenv("DATA_BASE_URL")
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}

	store := generated.NewStore(conn)
	customQuerier := NewCustomQuerier(conn)
	userRepo := repository4.NewUserRepository(ctx, store)
	saleTypeRepo := repository.NewSaleTypeRepository(ctx, store)
	saleManagerRepo := repository2.NewSalesManagerRepository(saleTypeRepo, ctx, store, customQuerier)
	branchRepo := repository3.NewBranchRepository(ctx, store, customQuerier, saleTypeRepo)
	directorRepo := repo.NewBranchDirectorRepository(ctx, store)
	rankingsRepo := NewRankingsRepository(ctx, customQuerier, branchRepo)
	salesManagerStatisticRepo := repository2.NewSalesManagerStatisticRepository(saleTypeRepo, ctx, store)
	postRepo := NewPostRepository(ctx, store)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthorizationService(ctx, userRepo)
	salesManagerService := service2.NewSalesManagerService(saleManagerRepo, salesManagerStatisticRepo, saleTypeRepo)
	branchService := service3.NewBranchService(branchRepo)
	saleTypeService := service4.NewSaleTypeService(saleTypeRepo)
	directorService := service5.NewBranchDirectorService(directorRepo)
	analyticsService := service6.NewAnalyticsService(rankingsRepo)

	encKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY")

	tokenService := service.NewTokenService(&encKey)
	server.userService = userService
	server.salesManagerService = salesManagerService
	server.saleTypeService = saleTypeService
	server.branchService = branchService
	server.tokenService = tokenService
	server.authService = authService
	server.directorService = directorService
	server.analyticsService = analyticsService
	server.postRepository = postRepo
}

// Start runs the HTTP server a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
