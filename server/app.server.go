package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/snowlynxsoftware/oto-api/config"
	"github.com/snowlynxsoftware/oto-api/server/controllers"
	"github.com/snowlynxsoftware/oto-api/server/database"
	"github.com/snowlynxsoftware/oto-api/server/database/repositories"
	"github.com/snowlynxsoftware/oto-api/server/middleware"
	"github.com/snowlynxsoftware/oto-api/server/services"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type AppServer struct {
	appConfig config.IAppConfig
	router    *chi.Mux
	dB        *database.AppDataSource
}

func NewAppServer(config config.IAppConfig) *AppServer {

	r := chi.NewRouter()
	r.Use(mid.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.GetCorsAllowedOrigin()},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	return &AppServer{
		appConfig: config,
		router:    r,
	}
}

func (s *AppServer) Start() {

	// Check if the app is running in production mode
	var isProductionMode = s.appConfig.GetCloudEnv() != "local"

	// Setup logger
	util.SetupZeroLogger(s.appConfig.IsDebugMode())

	// Connect to DB
	s.dB = database.NewAppDataSource()
	s.dB.Connect(s.appConfig.GetDBConnectionString())

	// Configure Repositories
	waitlistRepository := repositories.NewWaitlistRepository(s.dB)
	userRepository := repositories.NewUserRepository(s.dB)
	triviaRepository := repositories.NewTriviaRepository(s.dB)

	// Configure Services
	emailService := services.NewEmailService(s.appConfig.GetSendgridAPIKey(), services.NewEmailTemplates())
	cryptoService := services.NewCryptoService(s.appConfig.GetAuthHashPepper())
	tokenService := services.NewTokenService(s.appConfig.GetJWTSecretKey())
	authService := services.NewAuthService(userRepository, tokenService, cryptoService, emailService)
	triviaService := services.NewTriviaService(triviaRepository)
	waitlistService := services.NewWaitlistService(waitlistRepository)
	userService := services.NewUserService(userRepository)

	// Configure Middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepository, tokenService)

	// Configure Controllers
	s.router.Mount("/health", controllers.NewHealthController().MapController())
	s.router.Mount("/auth", controllers.NewAuthController(authMiddleware, authService, isProductionMode, s.appConfig.GetCookieDomain()).MapController())
	s.router.Mount("/trivia", controllers.NewTriviaController(triviaService, authMiddleware).MapController())
	s.router.Mount("/waitlist", controllers.NewWaitlistController(waitlistService).MapController())
	s.router.Mount("/users", controllers.NewUserController(userService, authMiddleware).MapController())

	util.LogInfo("Starting server on localhost:3000")
	log.Fatal(http.ListenAndServe("localhost:3000", s.router))
}
