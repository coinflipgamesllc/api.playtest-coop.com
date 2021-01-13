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
			auth.POST("/reset-password", authController.RequestResetPassword)
			auth.GET("/reset-password/:otp", authController.ResetPassword)

			auth.POST("/signup", authController.Signup)
			auth.POST("/login", authController.Login)
			auth.GET("/logout", authController.Logout)
			auth.GET("/verify-email/:id", authController.VerifyEmail)
		}

		eventController := container.EventController()
		events := v1.Group("/events")
		{
			events.GET("", eventController.ListEvents)
			events.POST("", container.Authenticated(), eventController.CreateEvent)
			events.GET("/:id", eventController.GetEvent)
			events.PUT("/:id", container.Authenticated(), eventController.UpdateEvent)
		}

		fileController := container.FileController()
		files := v1.Group("/files")
		{
			files.GET("/sign", container.Authenticated(), fileController.PresignUpload)
			files.POST("", container.Authenticated(), fileController.CreateFile)
			files.GET("", container.Authenticated(), fileController.ListUserFiles)
			files.PUT("/:id", container.Authenticated(), fileController.UpdateFile)
			files.DELETE("/:id", container.Authenticated(), fileController.DeleteFile)
		}

		gameController := container.GameController()
		games := v1.Group("/games")
		{
			games.GET("", gameController.ListGames)
			games.POST("", container.Authenticated(), gameController.CreateGame)
			games.GET("/:id", gameController.GetGame)
			games.PUT("/:id", container.Authenticated(), gameController.UpdateGame)

			games.GET("/:id/rules", gameController.GetRules)
		}

		v1.GET("/mechanics", gameController.ListAvailableMechanics)

		userController := container.UserController()
		users := v1.Group("/users")
		{
			users.GET("", container.Authenticated(), userController.ListUsers)
		}
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

func serverError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
