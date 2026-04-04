package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"track-selection/internal/application/student"
	"track-selection/internal/config"
	"track-selection/internal/domain/shared/auth"
	"track-selection/internal/domain/shared/events"
	"track-selection/internal/infrastructure/eventbus"
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
	AuthRepo *postgres.AuthRepository

	// Технические сервисы
	JwtService *auth.JWTService

	// БД пул
	DB *pgxpool.Pool
}

// UseCases — только бизнес-сценарии (координация)
type UseCases struct {
	// Auth Use Cases
	RegisterUC *auth.RegisterUseCase
	LoginUC    *auth.LoginUseCase

	// Student Use Cases
	SelectTrackUC        *student.SelectTrackUseCase
	UpdateProfileUC      *student.UpdateProfileUseCase
	GetRecommendationsUC *student.GetRecommendationsUseCase

	// Admin Use Cases
	CreateTrackUC *admin.CreateTrackUseCase
	UpdateTrackUC *admin.UpdateTrackUseCase
	DeleteTrackUC *admin.DeleteTrackUseCase
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
	slog.Info("PostgreSQL initialized")

	// 2. Репозитории
	a.infra.AuthRepo = postgres.NewAuthRepository(poolPG)
	a.infra.StudentRepo = postgres.NewStudentRepository(poolPG)
	a.infra.AdminRepo = postgres.NewAdminRepository(poolPG)
	a.infra.TrackRepo = postgres.NewTrackRepository(poolPG)
	a.infra.SelectionRepo = postgres.NewSelectionRepository(poolPG)
	slog.Info("Repositories initialized")

	// 3. JWT сервис
	if a.cfg.Jwt_secret_key == "" {
		slog.Error("JWT secret key is required")
		os.Exit(1)
	}
	a.infra.JwtService = jwt.NewJWTService(&jwt.JWTConfig{
		Secret:     a.cfg.Jwt_secret_key,
		Expiration: 24 * time.Hour,
	})
	slog.Info("JWT service initialized")
}

// initEventBus — инициализация шины событий
func (a *App) initEventBus() {
	a.eventBus = eventbus.NewMemoryBus()
	slog.Info("EventBus initialized")
}

// initUseCases — инициализация всех Use Case'ов
func (a *App) initUseCases() {
	// Auth Use Cases
	a.useCases.RegisterUC = auth.NewRegisterUseCase(
		a.infra.AuthRepo,
		a.infra.StudentRepo,
		a.infra.AdminRepo,
		a.eventBus,
	)

	a.useCases.LoginUC = auth.NewLoginUseCase(
		a.infra.AuthRepo,
		a.infra.JwtService,
	)

	// Student Use Cases
	a.useCases.SelectTrackUC = student.NewSelectTrackUseCase(
		a.infra.StudentRepo,
		a.infra.TrackRepo,
		a.infra.SelectionRepo,
		a.eventBus,
	)

	a.useCases.UpdateProfileUC = student.NewUpdateProfileUseCase(
		a.infra.StudentRepo,
		a.eventBus,
	)

	a.useCases.GetRecommendationsUC = student.NewGetRecommendationsUseCase(
		a.infra.StudentRepo,
		a.infra.TrackRepo,
		a.infra.SelectionRepo,
	)

	// Admin Use Cases
	a.useCases.CreateTrackUC = admin.NewCreateTrackUseCase(
		a.infra.TrackRepo,
		a.eventBus,
	)

	a.useCases.UpdateTrackUC = admin.NewUpdateTrackUseCase(
		a.infra.TrackRepo,
		a.eventBus,
	)

	a.useCases.DeleteTrackUC = admin.NewDeleteTrackUseCase(
		a.infra.TrackRepo,
		a.eventBus,
	)

	slog.Info("Use Cases initialized")
}

// initHTTP — инициализация HTTP сервера и маршрутов
func (a *App) initHTTP() {
	handler := handlers.NewHandler(
		a.infra.JwtService,
		a.useCases.RegisterUC,
		a.useCases.LoginUC,
		a.useCases.SelectTrackUC,
		a.useCases.UpdateProfileUC,
		a.useCases.GetRecommendationsUC,
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
	r.HandleFunc("/_info", handler.TestHandler).Methods(http.MethodGet)
	r.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)

	// Эндпоинты с авторизацией (любая роль)
	// r.HandleFunc("/student/profile",
	// 	middleware.WithAuth(a.infra.JwtService, handler.GetProfile, middleware.RoleAny)).
	// 	Methods(http.MethodGet)

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
