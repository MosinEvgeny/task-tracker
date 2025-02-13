package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MosinEvgeny/task-tracker/internal/config"
	"github.com/MosinEvgeny/task-tracker/internal/handlers"
	"github.com/MosinEvgeny/task-tracker/internal/repository/postgres"
	"github.com/MosinEvgeny/task-tracker/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	config config.Config
	router *mux.Router
	db     *postgres.PostgresDB
}

func NewApp(config config.Config) (*App, error) {
	db, err := postgres.NewPostgresDB(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &App{
		config: config,
		router: mux.NewRouter(),
		db:     db,
	}, nil
}

func (a *App) Run() error {
	// Инициализация зависимостей
	userRepo := postgres.NewUserRepository(a.db)
	userService := service.NewUserService(userRepo)

	refreshTokenRepo := postgres.NewRefreshTokenRepository(a.db)
	refreshTokenService := service.NewRefreshTokenService(refreshTokenRepo)

	userHandler := handlers.NewUserHandler(userService, refreshTokenService, a.config)

	taskRepo := postgres.NewTaskRepository(a.db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	labelRepo := postgres.NewLabelRepository(a.db)
	labelService := service.NewLabelService(labelRepo)
	labelHandler := handlers.NewLabelHandler(labelService)

	// Настройка middleware
	authMiddleware := handlers.NewAuthMiddleware(userService, a.config)
	logMiddleware := handlers.Log

	// Настройка маршрутов
	a.router.HandleFunc("/register", userHandler.RegisterUser).Methods("POST")
	a.router.HandleFunc("/login", userHandler.LoginUser).Methods("POST")
	a.router.HandleFunc("/refresh", userHandler.RefreshToken).Methods("POST")

	// Маршрут для отзыва всех refresh токенов
	userRouter := a.router.PathPrefix("/users").Subrouter()
	userRouter.Use(authMiddleware.Authenticate)
	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")
	userRouter.HandleFunc("/revoke", userHandler.RevokeAllRefreshTokens).Methods("POST")

	taskRouter := a.router.PathPrefix("/tasks").Subrouter()
	taskRouter.Use(authMiddleware.Authenticate)
	taskRouter.HandleFunc("", taskHandler.CreateTask).Methods("POST")
	taskRouter.HandleFunc("/{id}", taskHandler.GetTask).Methods("GET")
	taskRouter.HandleFunc("/{id}", taskHandler.UpdateTask).Methods("PUT")
	taskRouter.HandleFunc("/{id}", taskHandler.DeleteTask).Methods("DELETE")

	labelRouter := a.router.PathPrefix("/labels").Subrouter()
	labelRouter.Use(authMiddleware.Authenticate)
	labelRouter.HandleFunc("", labelHandler.CreateLabel).Methods("POST")
	labelRouter.HandleFunc("/{id}", labelHandler.GetLabel).Methods("GET")
	labelRouter.HandleFunc("/{id}", labelHandler.UpdateLabel).Methods("PUT")
	labelRouter.HandleFunc("/{id}", labelHandler.DeleteLabel).Methods("DELETE")

	// Логирование всех запросов
	a.router.Use(logMiddleware)

	// CORS настройки
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// CORS middleware в router
	handler := c.Handler(a.router)

	// Запуск сервера
	server := &http.Server{
		Addr:         ":" + a.config.AppPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 6. Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		log.Println("Server gracefully stopped")

		if err := a.db.Close(); err != nil {
			log.Fatalf("Database connection close failed: %v", err)
		}
		log.Println("Database connection closed")
	}()

	log.Printf("Starting server on port %s", a.config.AppPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
