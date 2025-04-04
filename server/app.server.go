package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"

	"github.com/snowlynxsoftware/oto-api/config"
	"github.com/snowlynxsoftware/oto-api/server/controllers"
	"github.com/snowlynxsoftware/oto-api/server/database"
	"github.com/snowlynxsoftware/oto-api/server/util"
)

type AppServer struct {
	AppConfig config.IAppConfig
	Router    *chi.Mux
	DB        *database.AppDataSource
}

func NewAppServer(config config.IAppConfig) *AppServer {

	r := chi.NewRouter()
	r.Use(mid.Logger)

	return &AppServer{
		AppConfig: config,
		Router:    r,
	}
}

func (s *AppServer) Start() {

	// Setup logger
	util.SetupZeroLogger(s.AppConfig.IsDebugMode())

	// Connect to DB
	s.DB = database.NewAppDataSource()
	s.DB.Connect(s.AppConfig.GetDBConnectionString())

	// Configure Repositories
	// userRepository := repositories.NewUserRepository(s.DB)
	// listRepository := repositories.NewListRepository(s.DB)

	// Configure Services
	// emailService := services.NewEmailService(s.AppConfig.SendgridAPIKey)
	// cryptoService := services.NewCryptoService(s.AppConfig.HashPepper)
	// tokenService := services.NewTokenService(s.AppConfig.JWTSecret)
	// authService := services.NewAuthService(userRepository, tokenService, cryptoService, emailService)

	// Configure Middleware
	//authMiddleware := middleware.NewAuthMiddleware(userService)

	// Configure Controllers
	s.Router.Mount("/health", controllers.NewHealthController().MapController())
	// s.Router.Mount("/api/settings", controllers.NewSettingsController().MapController())
	// s.Router.Mount("/api/users", controllers.NewUserController(userService).MapController())
	// s.Router.Mount("/api/auth", controllers.NewAuthController(authMiddleware, userService, authService).MapController())
	// s.Router.Mount("/api/lists", controllers.NewListController(authMiddleware, listService).MapController())
	// s.Router.Mount("/", controllers.NewUIController(&templatesFS, authMiddleware, listService).MapController())

	util.LogInfo(fmt.Sprintf("Starting server on port %s", "3000"))
	log.Fatal(http.ListenAndServe(":3000", s.Router))
}
