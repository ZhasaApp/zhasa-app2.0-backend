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

	router.POST("/images/upload", server.HandleUpload)

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
		adminRoute.GET("/branch/list", server.getBranches)
	}

	smRoute := router.Group("sales-manager/")
	smRoute.Use(getSalesManager(server.tokenService, server.salesManagerService))
	{
		smRoute.POST("/sale/new", server.saveSale)
		smRoute.GET("/branch/list", server.getBranches)
		smRoute.GET("/year-statistic", server.getYearStatistic)
	}

	branchRoute := router.Group("branch/").Use(verifyToken(server.tokenService))
	{
		branchRoute.GET("/year-statistic", server.getBranchYearStatistic)
	}

	directorRoute := router.Group("branch-director/")
	directorRoute.Use(getBranchDirector(server.tokenService, server.directorService))
	{
		directorRoute.POST("/goal/new", server.createSaleGoalForSalesManager)
	}

	router.GET("sales-manager/dashboard", server.getSalesManagerDashboardStatistic).Use(verifyToken(server.tokenService))

	router.GET("branch/dashboard", server.getBranchDashboardStatistic).Use(verifyToken(server.tokenService))

	router.GET("ratings/managers", server.getRankedSalesManager).Use(verifyToken(server.tokenService))

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

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthorizationService(ctx, userRepo)
	salesManagerService := service2.NewSalesManagerService(saleManagerRepo, saleTypeRepo)
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
}

// Start runs the HTTP server a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
