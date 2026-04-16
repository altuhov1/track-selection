package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"track-selection/internal/application/auth"
	authApp "track-selection/internal/application/auth"
	studApp "track-selection/internal/application/student"
	"track-selection/internal/application/track"
	"track-selection/internal/config"
	authDomain "track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/domain/student"
	"track-selection/internal/infrastructure/eventbus"
	authEventBus "track-selection/internal/infrastructure/eventbus/subscribers/auth"
	"track-selection/internal/infrastructure/http/handlers"
	"track-selection/internal/infrastructure/http/middleware"
	"track-selection/internal/infrastructure/jwt"
	"track-selection/internal/infrastructure/persistence/postgres"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg      *config.ConfigApp
	server   *http.Server
	useCases *UseCases
	infra    *Infrastructure
	eventBus events.EventBus
	rootCtx  context.Context
	cancel   context.CancelFunc
}

// Infrastructure — все технические компоненты (репозитории, сервисы, БД)
type Infrastructure struct {
	// Репозитории (работа с БД)
	AuthRepo              *postgres.AuthRepository
	StudentRepo           *postgres.StudentRepository
	AdminRepo             *postgres.AdminRepository
	PreferencesRepo       *postgres.PreferencesRepository
	ProfileCompletionRepo *postgres.ProfileCompletionRepository
	TrackRepo             *postgres.TrackRepository

	// Технические сервисы
	JwtService authDomain.JWTService

	// БД пул
	DB *pgxpool.Pool
}

// UseCases — только бизнес-сценарии (координация)
type UseCases struct {
	// Auth Use Cases
	RegisterUC *authApp.RegisterUseCase
	LoginUC    *authApp.LoginUseCase

	// Student Use Cases
	UpdatePreferencesUC    *studApp.UpdatePreferencesUseCase
	GetPreferencesUC       *studApp.GetPreferencesUseCase
	GetProfileCompletionUC *studApp.GetProfileCompletionUseCase

	GetAllTracksUC *track.GetAllTracksUseCase
	CreateTrackUC  *track.CreateTrackUseCase
	UpdateTrackUC  *track.UpdateTrackUseCase
	DeleteTrackUC  *track.DeleteTrackUseCase
}

func NewApp(cfg *config.ConfigApp) *App {
	rootCtx, cancel := context.WithCancel(context.Background())

	app := &App{
		cfg:      cfg,
		rootCtx:  rootCtx,
		cancel:   cancel,
		useCases: &UseCases{},
		infra:    &Infrastructure{},
	}

	app.initInfrastructure()
	app.initEventBus()
	app.initUseCases()
	app.initHTTP()

	return app
}

// initInfrastructure — инициализация всех технических компонентов
func (a *App) initInfrastructure() {
	// 1. Подключение к PostgreSQL
	poolPG, err := postgres.NewPoolPg(&postgres.PoolConfig{
		Host:     a.cfg.PG_DBHost,
		User:     a.cfg.PG_DBUser,
		Password: a.cfg.PG_DBPassword,
		DBName:   a.cfg.PG_DBName,
		SSLMode:  a.cfg.PG_DBSSLMode,
		Port:     a.cfg.PG_PORT,
	})
	if err != nil {
		slog.Error("Failed to initialize PostgreSQL", "error", err)
		os.Exit(1)
	}
	a.infra.DB = poolPG

	// 2. Репозитории
	a.infra.AuthRepo = postgres.NewAuthRepository(poolPG)
	a.infra.StudentRepo = postgres.NewStudentRepository(poolPG)
	a.infra.AdminRepo = postgres.NewAdminRepository(poolPG)
	a.infra.PreferencesRepo = postgres.NewPreferencesRepository(poolPG)
	a.infra.ProfileCompletionRepo = postgres.NewProfileCompletionRepository(poolPG)
	a.infra.TrackRepo = postgres.NewTrackRepository(poolPG)

	// Создаем дефолтные треки
	postgres.SeedTracks(a.rootCtx, a.infra.TrackRepo)

	// 3. JWT сервис
	if a.cfg.Jwt_secret_key == "" {
		slog.Error("JWT secret key is required")
		os.Exit(1)
	}
	a.infra.JwtService = jwt.NewJWTService(&jwt.JWTConfig{
		Secret:     a.cfg.Jwt_secret_key,
		Expiration: 24 * time.Hour,
	})
}

// initEventBus — инициализация шины событий
func (a *App) initEventBus() {
	a.eventBus = eventbus.NewMemoryBus()

	// Подписываем обработчики
	createStudentHandler := authEventBus.NewCreateStudentRegHandler(a.infra.StudentRepo)
	a.eventBus.Subscribe("student.registered", createStudentHandler)

	createAdminHandler := authEventBus.NewCreateAdminRegHandler(a.infra.AdminRepo)
	a.eventBus.Subscribe("admin.registered", createAdminHandler)

	createProfileHandler := authEventBus.NewCreateProfileCompletionHandler(
		a.infra.PreferencesRepo,
		a.infra.ProfileCompletionRepo,
	)
	a.eventBus.Subscribe("student.registered", createProfileHandler)
}

// initUseCases — инициализация всех Use Case'ов
func (a *App) initUseCases() {
	// Profile checker
	profileChecker := student.NewProfileChecker()

	// Auth Use Cases
	a.useCases.RegisterUC = auth.NewRegisterUseCase(
		a.infra.AuthRepo,
		a.eventBus,
	)

	a.useCases.LoginUC = auth.NewLoginUseCase(
		a.infra.AuthRepo,
		a.infra.JwtService,
	)

	// Student Use Cases
	a.useCases.UpdatePreferencesUC = studApp.NewUpdatePreferencesUseCase(
		a.infra.PreferencesRepo,
		a.infra.ProfileCompletionRepo,
		profileChecker,
		a.eventBus,
	)

	a.useCases.GetPreferencesUC = studApp.NewGetPreferencesUseCase(a.infra.PreferencesRepo)
	a.useCases.GetProfileCompletionUC = studApp.NewGetProfileCompletionUseCase(a.infra.ProfileCompletionRepo)
	a.useCases.GetAllTracksUC = track.NewGetAllTracksUseCase(a.infra.TrackRepo)
	a.useCases.CreateTrackUC = track.NewCreateTrackUseCase(a.infra.TrackRepo)
	a.useCases.UpdateTrackUC = track.NewUpdateTrackUseCase(a.infra.TrackRepo)
	a.useCases.DeleteTrackUC = track.NewDeleteTrackUseCase(a.infra.TrackRepo)
	slog.Info("Use Cases initialized")
}

// initHTTP — инициализация HTTP сервера и маршрутов
func (a *App) initHTTP() {
	handler := handlers.NewHandler(
		a.useCases.RegisterUC,
		a.useCases.LoginUC,
		a.useCases.UpdatePreferencesUC,
		a.useCases.GetPreferencesUC,
		a.useCases.GetProfileCompletionUC,
		a.useCases.GetAllTracksUC,
		a.useCases.CreateTrackUC,
		a.useCases.UpdateTrackUC,
		a.useCases.DeleteTrackUC,
	)

	router := a.setupRoutes(handler)

	a.server = &http.Server{
		Addr:         ":" + a.cfg.App_port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

// setupRoutes — настройка маршрутов
func (a *App) setupRoutes(handler *handlers.Handler) http.Handler {
	r := mux.NewRouter()

	// Публичные эндпоинты (без аутентификации)
	r.HandleFunc("/api/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/login", handler.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/all-tracks", handler.GetAllTracks).Methods(http.MethodGet)

	// Защищенные эндпоинты (с аутентификацией)
	r.HandleFunc("/api/me", middleware.WithAuth(a.infra.JwtService, handler.GetMe, middleware.RoleAny)).Methods(http.MethodGet)
	r.HandleFunc("/api/me/edit-info", middleware.WithAuth(a.infra.JwtService, handler.UpdatePreferences, middleware.RoleAny)).Methods(http.MethodPost)
	r.HandleFunc("/api/me/info", middleware.WithAuth(a.infra.JwtService, handler.GetPreferences, middleware.RoleAny)).Methods(http.MethodGet)
	r.HandleFunc("/api/me/profile-completion", middleware.WithAuth(a.infra.JwtService, handler.GetProfileCompletion, middleware.RoleAny)).Methods(http.MethodGet)
	r.HandleFunc("/api/new-track", middleware.WithAuth(a.infra.JwtService, handler.CreateTrack, middleware.RoleAdmin)).Methods(http.MethodPost)
	r.HandleFunc("/api/edit-track/{id}", middleware.WithAuth(a.infra.JwtService, handler.UpdateTrack, middleware.RoleAdmin)).Methods(http.MethodPut)
	r.HandleFunc("/api/delete-track/{id}", middleware.WithAuth(a.infra.JwtService, handler.DeleteTrack, middleware.RoleAdmin)).Methods(http.MethodDelete)

	return middleware.ContextMiddleware(a.rootCtx, r)
}

// Run — запуск сервера
func (a *App) Run() {
	go a.startServer()
	a.waitForShutdown()
}

func (a *App) startServer() {
	slog.Info("Server starting", "port", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func (a *App) waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("Shutting down gracefully...")

	a.cancel()
	a.shutdown()

	slog.Info("Application stopped")
}

func (a *App) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// Закрываем соединение с БД
	if a.infra != nil && a.infra.DB != nil {
		a.infra.DB.Close()
		slog.Info("Database connection closed")
	}
}
