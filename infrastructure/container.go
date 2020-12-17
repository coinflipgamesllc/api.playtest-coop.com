package infrastructure

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/app"
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/persistence"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/controller"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/events"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
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

	// Domain
	fileRepository domain.FileRepository
	gameRepository domain.GameRepository
	userRepository domain.UserRepository

	// Infrastructure
	db        *gorm.DB
	logger    *zap.Logger
	mail      mailgun.Mailgun
	router    *gin.Engine
	s3Bucket  string
	s3Client  *minio.Client
	session   sessions.Store
	templates map[string]*template.Template

	// UI
	authController *controller.AuthController
	fileController *controller.FileController
	gameController *controller.GameController

	authenticated gin.HandlerFunc

	eventHandler *events.EventHandler
}

// AuthService for handling authentication & authorization
func (c *Container) AuthService() *app.AuthService {
	if c.authService == nil {
		c.authService = &app.AuthService{
			AuthToken:      viper.GetString("AUTH_TOKEN"),
			Logger:         c.Logger(),
			UserRepository: c.UserRepository(),
		}
	}

	return c.authService
}

func (c *Container) FileService() *app.FileService {
	if c.fileService == nil {
		c.fileService = &app.FileService{
			FileRepository: c.FileRepository(),
			GameRepository: c.GameRepository(),
			UserRepository: c.UserRepository(),
			Logger:         c.Logger(),
			S3Bucket:       c.S3Bucket(),
			S3Client:       c.S3Client(),
		}
	}

	return c.fileService
}

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

func (c *Container) MailService() *app.MailService {
	if c.mailService == nil {
		c.mailService = &app.MailService{
			FromAddress: viper.GetString("FROM_ADDRESS"),
			Hostname:    viper.GetString("HOSTNAME"),
			MailClient:  c.Mail(),
			Templates:   c.Templates(),
		}
	}

	return c.mailService
}

func (c *Container) FileRepository() domain.FileRepository {
	if c.fileRepository == nil {
		c.fileRepository = &persistence.FileRepository{
			DB: c.DB(),
		}
	}

	return c.fileRepository
}

func (c *Container) GameRepository() domain.GameRepository {
	if c.gameRepository == nil {
		c.gameRepository = &persistence.GameRepository{
			DB: c.DB(),
		}
	}

	return c.gameRepository
}

func (c *Container) UserRepository() domain.UserRepository {
	if c.userRepository == nil {
		c.userRepository = &persistence.UserRepository{
			DB: c.DB(),
		}
	}

	return c.userRepository
}

func (c *Container) DB() *gorm.DB {
	if c.db == nil {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			viper.GetString("DB_HOSTNAME"),
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_DATABASE"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_SSLMODE"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		db.AutoMigrate(
			&domain.File{},
			&domain.Game{},
			&domain.User{},
		)

		c.db = db
	}

	return c.db
}

func (c *Container) Logger() *zap.Logger {
	if c.logger == nil {
		var logger *zap.Logger
		var err error
		if viper.GetString("ENVIRONMENT") == "development" {
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

func (c *Container) Mail() mailgun.Mailgun {
	if c.mail == nil {
		c.mail = mailgun.NewMailgun(
			viper.GetString("MAILGUN_DOMAIN"),
			viper.GetString("MAILGUN_APIKEY"),
		)
	}

	return c.mail
}

func (c *Container) Router() *gin.Engine {
	if c.router == nil {
		if viper.GetString("ENVIRONMENT") == "development" {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		c.router = gin.New()

		c.router.Use(sessions.Sessions("ptc_sess", c.Session()))
		c.router.Use(ginzap.Ginzap(c.Logger(), time.RFC3339, true))
		c.router.Use(ginzap.RecoveryWithZap(c.Logger(), true))
		c.router.LoadHTMLGlob("ui/template/error/*")
	}

	return c.router
}

func (c *Container) S3Bucket() string {
	if c.s3Bucket == "" {
		c.s3Bucket = viper.GetString("S3_BUCKET")
	}

	return c.s3Bucket
}

func (c *Container) S3Client() *minio.Client {
	if c.s3Client == nil {
		s3, err := minio.New(viper.GetString("S3_ENDPOINT"), &minio.Options{
			Creds:  credentials.NewStaticV4(viper.GetString("AWS_ACCESS_KEY"), viper.GetString("AWS_ACCESS_SECRET"), ""),
			Secure: true,
		})

		if err != nil {
			log.Fatal(err)
		}

		c.s3Client = s3
	}

	return c.s3Client
}

func (c *Container) Session() sessions.Store {
	if c.session == nil {
		c.session = cookie.NewStore([]byte(viper.GetString("AUTH_TOKEN")))

		secure := true
		if viper.GetString("ENVIRONMENT") == "development" {
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

func (c *Container) FileController() *controller.FileController {
	if c.fileController == nil {
		c.fileController = &controller.FileController{
			FileService: c.FileService(),
		}
	}

	return c.fileController
}

func (c *Container) GameController() *controller.GameController {
	if c.gameController == nil {
		c.gameController = &controller.GameController{
			GameService: c.GameService(),
		}
	}

	return c.gameController
}

// Authenticated middleware for ensuring that an HTTP request includes a valid access token
func (c *Container) Authenticated() gin.HandlerFunc {
	if c.authenticated == nil {
		c.authenticated = middleware.Authenticated(viper.GetString("AUTH_TOKEN"))
	}

	return c.authenticated
}

func (c *Container) EventHandler() *events.EventHandler {
	if c.eventHandler == nil {
		c.eventHandler = &events.EventHandler{
			MailService: c.MailService(),
			Logger:      c.Logger(),
		}
	}

	return c.eventHandler
}
