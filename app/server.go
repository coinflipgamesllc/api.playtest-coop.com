package app

import (
	"fmt"
	"html/template"
	"log"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	authToken string
	hostname  string

	mail          mailgun.Mailgun
	mailValidator mailgun.EmailValidator

	router *gin.Engine

	s3Bucket string
	s3Client *minio.Client

	templates map[string]*template.Template

	fileRepository domain.FileRepository
	gameRepository domain.GameRepository
	userRepository domain.UserRepository
}

func NewServer() *Server {
	db := db()

	// Create our server
	server := &Server{
		authToken:      viper.GetString("AUTH_TOKEN"),
		hostname:       viper.GetString("HOSTNAME"),
		mail:           mail(),
		router:         gin.Default(),
		s3Bucket:       viper.GetString("S3_BUCKET"),
		s3Client:       s3(),
		templates:      templates(),
		fileRepository: &persistence.FileRepository{DB: db},
		gameRepository: &persistence.GameRepository{DB: db},
		userRepository: &persistence.UserRepository{DB: db},
	}

	// Register routes
	server.routes()
	server.router.LoadHTMLGlob("ui/template/error/*")

	return server
}

func (s *Server) Run() {
	// Start events handlers
	go func() {
		s.listenForEvents()
	}()

	// Start http handler
	s.router.Run(":" + viper.GetString("PORT"))
}

func db() *gorm.DB {
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

	return db
}

func mail() mailgun.Mailgun {
	return mailgun.NewMailgun(
		viper.GetString("MAILGUN_DOMAIN"),
		viper.GetString("MAILGUN_APIKEY"),
	)
}

func s3() *minio.Client {
	s3, err := minio.New(viper.GetString("S3_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("AWS_ACCESS_KEY"), viper.GetString("AWS_ACCESS_SECRET"), ""),
		Secure: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	return s3
}

func templates() map[string]*template.Template {
	t := map[string]*template.Template{}

	basePath := "ui/template/"
	paths := []string{
		"email/verify-email",
		"email/welcome",
		"error/404",
		"error/500",
	}
	for _, p := range paths {
		tpl, err := template.ParseFiles(basePath + p + ".html")
		if err != nil {
			log.Fatal(err)
		}

		t[p] = tpl
	}

	return t
}
