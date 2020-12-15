package app

// type Server struct {
// 	authToken string
// 	hostname  string

// 	fileRepository domain.FileRepository
// 	gameRepository domain.GameRepository
// 	userRepository domain.UserRepository
// }

// func NewServer() *Server {
// 	// TODO - move to container in infrastructure
// 	db := db()

// 	logger := zap.S()

// 	fileRepository := &persistence.FileRepository{DB: db}
// 	gameRepository := &persistence.GameRepository{DB: db}
// 	userRepository := &persistence.UserRepository{DB: db}

// 	authService := &AuthService{
// 		AuthToken:      viper.GetString("AUTH_TOKEN"),
// 		Logger:         logger,
// 		UserRepository: userRepository,
// 	}

// 	authController := &controller.AuthController{}

// 	// Create our server
// 	server := &Server{
// 		authToken:      viper.GetString("AUTH_TOKEN"),
// 		hostname:       viper.GetString("HOSTNAME"),
// 		mail:           mail(),
// 		router:         gin.Default(),
// 		s3Bucket:       viper.GetString("S3_BUCKET"),
// 		s3Client:       s3(),
// 		templates:      templates(),
// 		authService:    authService,
// 		authController: authController,
// 		fileRepository: fileRepository,
// 		gameRepository: gameRepository,
// 		userRepository: userRepository,
// 	}

// 	// Register routes
// 	server.routes()
// 	server.router.LoadHTMLGlob("ui/template/error/*")

// 	return server
// }

// func (s *Server) Run() {
// 	// Start events handlers
// 	go func() {
// 		s.listenForEvents()
// 	}()

// 	// Start http handler
// 	s.router.Run(":" + viper.GetString("PORT"))
// }
