package ui

import (
	_ "github.com/coinflipgamesllc/api.playtest-coop.com/docs" // Required to include swagger docs

	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes adds all the api routes to the application
func RegisterRoutes(container *infrastructure.Container) {
	router := container.Router()

	v1 := router.Group("/v1")
	{
		authController := container.AuthController()
		auth := v1.Group("/auth")
		{
			auth.GET("/user", container.Authenticated(), authController.GetUser)
			auth.PUT("/user", container.Authenticated(), authController.UpdateUser)

			auth.POST("/signup", authController.Signup)
			auth.POST("/login", authController.Login)
			auth.POST("/token", authController.RefreshToken)
			auth.GET("/verify-email/:id", authController.VerifyEmail)
		}

		fileController := container.FileController()
		files := v1.Group("/files")
		{
			files.GET("/sign", container.Authenticated(), fileController.PresignUpload)
			files.POST("", container.Authenticated(), fileController.CreateFile)
			files.GET("", container.Authenticated(), fileController.ListUserFiles)
			files.DELETE("/:id", container.Authenticated(), fileController.DeleteFile)
		}

		// games := v1.Group("/games")
		// {
		// 	games.GET("", s.handleListGames())
		// 	games.POST("", container.Authenticated(), s.handleCreateGame())
		// 	games.GET("/:id", s.handleGetGame())
		// 	games.PUT("/:id", container.Authenticated(), s.handleUpdateGame())
		// }
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func serverError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
