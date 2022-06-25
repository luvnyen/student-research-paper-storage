package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luvnyen/student-research-paper-storage/cmd/config"

	"github.com/luvnyen/student-research-paper-storage/middleware"
	controller "github.com/luvnyen/student-research-paper-storage/pkg/controllers"
	"github.com/luvnyen/student-research-paper-storage/repository"
	"github.com/luvnyen/student-research-paper-storage/service"
	"gorm.io/gorm"
)

var (
	db                *gorm.DB                     = config.SetupDatabaseConnection()
	studentRepository repository.StudentRepository = repository.NewStudentRepository(db)
	paperRepository   repository.PaperRepository   = repository.NewPaperRepository(db)
	jwtService        service.JWTService           = service.NewJWTService()
	authService       service.AuthService          = service.NewAuthService(studentRepository)
	paperService      service.PaperService         = service.NewPaperService(paperRepository)
	authController    controller.AuthController    = controller.NewAuthController(authService, jwtService)
	paperController   controller.PaperController   = controller.NewPaperController(paperService, jwtService)
)

func main() {
	router := gin.Default()

	authRoutes := router.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
		authRoutes.POST("/register", authController.Register)
	}

	paperRoutes := router.Group("api/paper", middleware.AuthorizeJWT(jwtService))
	{
		paperRoutes.GET("/", paperController.All)
		paperRoutes.GET("/:id", paperController.FindByID)
		paperRoutes.POST("/", paperController.InsertPaper)
		paperRoutes.POST("/search", paperController.FindByTitleAuthorAbstract)
		paperRoutes.GET("/download/:id", paperController.DownloadPaper)
	}

	router.Run("127.0.0.1:3000")
}
