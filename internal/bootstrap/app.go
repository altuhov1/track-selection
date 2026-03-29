package bootstrap

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"track-selection/internal/config"
	"track-selection/internal/infrastructure/http/handlers"

	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg      *config.ConfigApp
	server   *http.Server
	useCases *UseCases     // вместо services
	repos    *Repositories // вместо storages
	eventBus domain.EventBus
	rootCtx  context.Context
	cancel   context.CancelFunc
}

// Repositories — теперь это интерфейсы, а не конкретные реализации
// (хотя реализации всё равно postgres)
type Repositories struct {
	Room     domain.RoomRepository
	Schedule domain.ScheduleRepository
	Slot     domain.SlotRepository
	Booking  domain.BookingRepository
	User     domain.UserRepository
	pool     *pgxpool.Pool
}

// UseCases — тонкий слой оркестрации (бывшие services, но тоньше)
type UseCases struct {
	Room     *room.UseCases
	Schedule *schedule.UseCases
	Slot     *slot.UseCases
	Booking  *booking.UseCases
	Auth     *auth.UseCases
}

func NewApp(cfg *config.ConfigApp) *App {
	rootCtx, cancel := context.WithCancel(context.Background())

	app := &App{
		cfg:     cfg,
		rootCtx: rootCtx,
		cancel:  cancel,
	}

	app.initRepositories() // было initStorages
	app.initEventBus()     // НОВО: инициализация шины событий
	app.initUseCases()     // было initServices
	app.initHTTP()         // осталось без изменений

	// Запускаем обработчики событий (воркеры)
	app.startEventHandlers()

	// Запускаем фоновые задачи
	app.startBackgroundJobs()

	return app
}

func (a *App) initRepositories() {
	poolPG, err := postgres.NewPool(a.cfg) // было storage.NewPoolPg
	if err != nil {
		slog.Error("Failed to initialize PostgreSQL pool", "error", err)
		os.Exit(1)
	}

	// Теперь репозитории — это реализации интерфейсов
	a.repos = &Repositories{
		Room:     postgres.NewRoomRepository(poolPG),     // было storage.NewRoomStorage
		Schedule: postgres.NewScheduleRepository(poolPG), // было storage.NewScheduleStorage
		Slot:     postgres.NewSlotRepository(poolPG),     // было storage.NewSlotStorage
		Booking:  postgres.NewBookingRepository(poolPG),  // было storage.NewBookingStorage
		User:     postgres.NewUserRepository(poolPG),     // было storage.NewUserStorage
		pool:     poolPG,
	}
}

func (a *App) initEventBus() {
	// Выбираем реализацию в зависимости от конфига
	if a.cfg.UseKafka {
		// Kafka для продакшена
		a.eventBus = events.NewKafkaEventBus(a.cfg.KafkaBrokers)
	} else {
		// In-memory для разработки/тестов
		a.eventBus = events.NewInMemoryEventBus()
	}
}

func (a *App) initUseCases() {
	// Создаём JWT сервис (это чистая логика без зависимостей)
	jwtService := auth.NewJWTService(a.cfg.JwtSecretKey) // было services.NewService

	// Создаём генератор слотов (доменный сервис)
	slotGenerator := schedule.NewSlotGenerator(14) // было services.NewSlotGenerator

	// Use Cases для комнат
	roomUseCases := room.NewUseCases(
		a.repos.Room,
		a.repos.Booking, // нужно для проверок
		a.eventBus,
	)

	// Use Cases для расписаний
	scheduleUseCases := schedule.NewUseCases(
		a.repos.Room,
		a.repos.Schedule,
		a.repos.Slot,
		slotGenerator,
		a.eventBus,
	)

	// Use Cases для слотов
	slotUseCases := slot.NewUseCases(
		a.repos.Room,
		a.repos.Schedule,
		a.repos.Slot,
		a.repos.Booking,
		a.eventBus,
	)

	// Use Cases для бронирований
	bookingUseCases := booking.NewUseCases(
		a.repos.Slot,
		a.repos.Booking,
		a.eventBus,
	)

	// Use Cases для аутентификации
	authUseCases := auth.NewUseCases(
		a.repos.User,
		jwtService,
		a.eventBus,
	)

	a.useCases = &UseCases{
		Room:     roomUseCases,
		Schedule: scheduleUseCases,
		Slot:     slotUseCases,
		Booking:  bookingUseCases,
		Auth:     authUseCases,
	}
}

func (a *App) startEventHandlers() {
	// Подписываемся на события и запускаем обработчики
	// Например, когда создаётся бронирование — отправляем уведомление

	a.eventBus.Subscribe("booking.created", func(ctx context.Context, event domain.Event) error {
		// Здесь может быть отправка email, уведомление в Slack и т.д.
		slog.Info("Booking created event received", "booking_id", event.AggregateID())
		return nil
	})

	a.eventBus.Subscribe("slot.booked", func(ctx context.Context, event domain.Event) error {
		// Обновляем кэш, отправляем WebSocket уведомление
		return nil
	})
}

func (a *App) startBackgroundJobs() {
	// Запускаем воркер для outbox паттерна (если используем)
	if a.cfg.UseOutbox {
		outboxWorker := events.NewOutboxWorker(a.repos.pool, a.eventBus)
		go outboxWorker.Start(a.rootCtx)
	}

	// Запускаем периодическую очистку устаревших слотов
	// Теперь это responsibility use case, а не сервиса
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-a.rootCtx.Done():
				return
			case <-ticker.C:
				if err := a.useCases.Slot.CleanupExpiredSlots(a.rootCtx); err != nil {
					slog.Error("Failed to cleanup expired slots", "error", err)
				}
			}
		}
	}()
}

func (a *App) initHTTP() {
	// Создаём хендлеры с новыми use cases
	handler := handlers.NewHandler(
		a.useCases.Auth,
		a.useCases.Room,
		a.useCases.Schedule,
		a.useCases.Slot,
		a.useCases.Booking,
	)

	router := a.setupRoutes(handler)

	a.server = &http.Server{
		Addr:         ":" + a.cfg.AppPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func (a *App) setupRoutes(handler *handlers.Handler) http.Handler {
	r := mux.NewRouter()

	// Публичные эндпоинты
	r.HandleFunc("/_info", handler.TestHandler).Methods(http.MethodGet)

	// Эндпоинты с авторизацией (любая роль)
	// r.HandleFunc("/rooms/list", middleware.WithAuth(handler.ListRooms, middleware.RoleAny)).Methods(http.MethodGet)
	// r.HandleFunc("/rooms/{roomId}/slots/list", middleware.WithAuth(handler.ListSlots, middleware.RoleAny)).Methods(http.MethodGet)

	// Добавляем middleware с JWT сервисом (теперь он встроен в хендлеры)
	return middleware.ContextMiddleware(a.rootCtx, r)
}

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

	// Закрываем event bus
	if err := a.eventBus.Close(); err != nil {
		slog.Error("Failed to close event bus", "error", err)
	}

	a.shutdown()

	slog.Info("Application stopped")
}

func (a *App) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
}
