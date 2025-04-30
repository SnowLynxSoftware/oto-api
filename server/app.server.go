package server

import (
	"fmt"
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
		AllowedOrigins:   []string{"*"}, // Change this to your frontend's URL in production
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

	// Configure Middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepository, tokenService)

	// Configure Controllers
	s.router.Mount("/health", controllers.NewHealthController().MapController())
	s.router.Mount("/auth", controllers.NewAuthController(authMiddleware, authService).MapController())
	s.router.Mount("/trivia", controllers.NewTriviaController(triviaService, authMiddleware).MapController())
	s.router.Mount("/waitlist", controllers.NewWaitlistController(waitlistService).MapController())

	util.LogInfo(fmt.Sprintf("Starting server on port %s", "3000"))
	log.Fatal(http.ListenAndServe(":3000", s.router))
}
