package infrastructure

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/persistence"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/validation"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/controller"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/events"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container is a lazy-load dependency injection container
type Container struct {
	// Application
	authService *app.AuthService
	fileService *app.FileService
	gameService *app.GameService
	mailService *app.MailService
	userService *app.UserService

	// Domain
	fileRepository domain.FileRepository
	gameRepository domain.GameRepository
	userRepository domain.UserRepository

	// Infrastructure
	db        *gorm.DB
	logger    *zap.Logger
	mail      mailgun.Mailgun
	router    *gin.Engine
	s3Client  *minio.Client
	session   sessions.Store
	templates map[string]*template.Template

	// UI
	authController *controller.AuthController
	fileController *controller.FileController
	gameController *controller.GameController
	userController *controller.UserController

	authenticated gin.HandlerFunc

	eventHandler *events.EventHandler
}

// AuthService for handling authentication & authorization
func (c *Container) AuthService() *app.AuthService {
	if c.authService == nil {
		c.authService = &app.AuthService{
			AuthToken:      os.Getenv("AUTH_TOKEN"),
			Logger:         c.Logger(),
			UserRepository: c.UserRepository(),
		}
	}

	return c.authService
}

// FileService for handling file uploads/downloads/etc
func (c *Container) FileService() *app.FileService {
	if c.fileService == nil {
		c.fileService = &app.FileService{
			FileRepository: c.FileRepository(),
			GameRepository: c.GameRepository(),
			UserRepository: c.UserRepository(),
			Logger:         c.Logger(),
			S3Bucket:       os.Getenv("S3_BUCKET"),
			S3Client:       c.S3Client(),
		}
	}

	return c.fileService
}

// GameService for general game content interaction
func (c *Container) GameService() *app.GameService {
	if c.gameService == nil {
		c.gameService = &app.GameService{
			GameRepository: c.GameRepository(),
			UserRepository: c.UserRepository(),
			Logger:         c.Logger(),
		}
	}

	return c.gameService
}

// MailService handles emailing users
func (c *Container) MailService() *app.MailService {
	if c.mailService == nil {
		c.mailService = &app.MailService{
			FromAddress: os.Getenv("FROM_ADDRESS"),
			Hostname:    os.Getenv("HOSTNAME"),
			MailClient:  c.Mail(),
			Templates:   c.Templates(),
		}
	}

	return c.mailService
}

// UserService for general user content interaction
func (c *Container) UserService() *app.UserService {
	if c.userService == nil {
		c.userService = &app.UserService{
			UserRepository: c.UserRepository(),
			Logger:         c.Logger(),
		}
	}

	return c.userService
}

// FileRepository implementation for database
func (c *Container) FileRepository() domain.FileRepository {
	if c.fileRepository == nil {
		c.fileRepository = &persistence.FileRepository{
			DB: c.DB(),
		}
	}

	return c.fileRepository
}

// GameRepository implementation for database
func (c *Container) GameRepository() domain.GameRepository {
	if c.gameRepository == nil {
		c.gameRepository = &persistence.GameRepository{
			DB: c.DB(),
		}
	}

	return c.gameRepository
}

// UserRepository implementation for database
func (c *Container) UserRepository() domain.UserRepository {
	if c.userRepository == nil {
		c.userRepository = &persistence.UserRepository{
			DB: c.DB(),
		}
	}

	return c.userRepository
}

// DB adapter for postgresql
func (c *Container) DB() *gorm.DB {
	if c.db == nil {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			os.Getenv("DB_HOSTNAME"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_DATABASE"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSLMODE"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		db.AutoMigrate(
			&domain.File{},
			&domain.Game{},
			&domain.User{},
			&game.RulesSection{},
		)

		c.db = db
	}

	return c.db
}

// Logger for consistent logging, application-wide
func (c *Container) Logger() *zap.Logger {
	if c.logger == nil {
		var logger *zap.Logger
		var err error
		if os.Getenv("ENVIRONMENT") == "development" {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatal(err)
		}

		c.logger = logger
	}

	return c.logger
}

// Mail client for mailgun
func (c *Container) Mail() mailgun.Mailgun {
	if c.mail == nil {
		c.mail = mailgun.NewMailgun(
			os.Getenv("MAILGUN_DOMAIN"),
			os.Getenv("MAILGUN_APIKEY"),
		)
	}

	return c.mail
}

// Router sets up the gin router
func (c *Container) Router() *gin.Engine {
	if c.router == nil {
		if os.Getenv("ENVIRONMENT") == "development" {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		c.router = gin.New()

		// Register the validator error formatter for later
		validation.NewJSONFormatter()

		c.router.Use(sessions.Sessions("ptc_sess", c.Session()))
		c.router.Use(ginzap.Ginzap(c.Logger(), time.RFC3339, true))
		c.router.Use(ginzap.RecoveryWithZap(c.Logger(), true))
		c.router.LoadHTMLGlob("ui/template/error/*")
	}

	return c.router
}

// S3Client for talking to s3-compatible storage
func (c *Container) S3Client() *minio.Client {
	if c.s3Client == nil {
		s3, err := minio.New(os.Getenv("S3_ENDPOINT"), &minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_ACCESS_SECRET"), ""),
			Secure: true,
		})

		if err != nil {
			log.Fatal(err)
		}

		c.s3Client = s3
	}

	return c.s3Client
}

// Session storage
func (c *Container) Session() sessions.Store {
	if c.session == nil {
		c.session = cookie.NewStore([]byte(os.Getenv("AUTH_TOKEN")))

		secure := true
		if os.Getenv("ENVIRONMENT") == "development" {
			secure = false
		}
		c.session.Options(sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   60 * 60, // 1 Hour
			Secure:   secure,
			HttpOnly: true,
		})
	}

	return c.session
}

// Templates initializes all the templates used by the MailService
func (c *Container) Templates() map[string]*template.Template {
	if c.templates == nil {
		t := map[string]*template.Template{}

		basePath := "ui/template/"
		paths := []string{
			"email/reset-password",
			"email/verify-email",
			"email/welcome",
		}
		for _, p := range paths {
			tpl, err := template.ParseFiles(basePath + p + ".html")
			if err != nil {
				log.Fatal(err)
			}

			t[p] = tpl
		}

		c.templates = t
	}

	return c.templates
}

// AuthController for handling /auth routes
func (c *Container) AuthController() *controller.AuthController {
	if c.authController == nil {
		c.authController = &controller.AuthController{
			AuthService: c.AuthService(),
		}
	}

	return c.authController
}

// FileController for handling /files routes
func (c *Container) FileController() *controller.FileController {
	if c.fileController == nil {
		c.fileController = &controller.FileController{
			FileService: c.FileService(),
		}
	}

	return c.fileController
}

// GameController for handling /games routes
func (c *Container) GameController() *controller.GameController {
	if c.gameController == nil {
		c.gameController = &controller.GameController{
			GameService: c.GameService(),
		}
	}

	return c.gameController
}

// UserController for handling /users routes
func (c *Container) UserController() *controller.UserController {
	if c.userController == nil {
		c.userController = &controller.UserController{
			UserService: c.UserService(),
		}
	}

	return c.userController
}

// Authenticated middleware for ensuring that an HTTP request includes a valid access token
func (c *Container) Authenticated() gin.HandlerFunc {
	if c.authenticated == nil {
		c.authenticated = middleware.Authenticated(os.Getenv("AUTH_TOKEN"))
	}

	return c.authenticated
}

// EventHandler for handling domain events subscribers
func (c *Container) EventHandler() *events.EventHandler {
	if c.eventHandler == nil {
		c.eventHandler = &events.EventHandler{
			MailService: c.MailService(),
			Logger:      c.Logger(),
		}
	}

	return c.eventHandler
}
