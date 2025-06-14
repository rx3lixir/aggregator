package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rx3lixir/agg-api/internal/db"
	"github.com/rx3lixir/agg-api/internal/lib/logger"
	"github.com/rx3lixir/agg-api/token"
)

// APIServer представляет сервер API с настройками и обработчиками.
type APIServer struct {
	listenAddr string
	logger     logger.Logger
	server     *http.Server
	store      db.Storage
	sessions   db.SessionStorage
	dbContext  context.Context
	TokenMaker *token.JWTMaker
}

// NewAPIServer создает новый экземпляр APIServer с указанным адресом.
func NewAPIServer(listenAddr string, log logger.Logger, store db.Storage, sessions db.SessionStorage, dbContext context.Context, secretKey string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		logger:     log,
		store:      store,
		sessions:   sessions,
		dbContext:  dbContext,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()

	tokenMaker := s.TokenMaker

	// --- ПОЛЬЗОВАТЕЛИ --- \\
	router.Route("/user", func(r chi.Router) {
		// Публичные маршруты для пользователей
		r.Post("/", s.makeHTTPHandleFunc(s.handleCreateUser))
		r.Post("/login", s.makeHTTPHandleFunc(s.handleLoginUser))

		// Управление пользователями (только для админов)
		r.Group(func(r chi.Router) {
			r.Use(GetAdminAuthMiddleWareFunc(tokenMaker))
			r.Get("/", s.makeHTTPHandleFunc(s.handleGetUsers))
			r.Delete("/{id}", s.makeHTTPHandleFunc(s.handleDeleteUser))
		})

		// Действия залогиненного пользователя
		r.Group(func(r chi.Router) {
			r.Use(GetAuthMiddleWareFunc(tokenMaker))

			r.Route("/{id}", func(r chi.Router) {
				r.Put("/", s.makeHTTPHandleFunc(s.handleUpdateUser))
				r.Get("/", s.makeHTTPHandleFunc(s.handleGetUserById))
			})
			r.Post("/logout", s.makeHTTPHandleFunc(s.handleLogoutUser))
		})

		// Управление токенами (только для админов)
		r.Route("/tokens", func(r chi.Router) {
			r.Use(GetAdminAuthMiddleWareFunc(tokenMaker))
			r.Post("/renew", s.makeHTTPHandleFunc(s.handleRenewAcessToken))
			r.Post("/revoke", s.makeHTTPHandleFunc(s.handleRevokeSession))
		})
	})

	// --- СОБЫТИЯ --- \\
	router.Route("/events", func(r chi.Router) {
		// Публичные маршруты для событий
		r.Get("/", s.makeHTTPHandleFunc(s.handleGetEvents))
		r.Get("/{id}", s.makeHTTPHandleFunc(s.handleGetEventById))
		r.Get("/category/{categoryId}", s.makeHTTPHandleFunc(s.handleGetEventsByCategory))

		// Защищенные маршруты для событий (только для админов)
		r.Group(func(r chi.Router) {
			r.Use(GetAdminAuthMiddleWareFunc(tokenMaker))
			r.Post("/", s.makeHTTPHandleFunc(s.handleCreateEvent))

			r.Route("/{id}", func(r chi.Router) {
				r.Put("/", s.makeHTTPHandleFunc(s.handleUpdateEvent))
				r.Delete("/", s.makeHTTPHandleFunc(s.handleDeleteEvent))
			})
		})
	})

	s.server = &http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	// Настраиваем корректное завершение работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		s.logger.Info("API server starting", "address", s.listenAddr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Fatal error starting server", "error", err)
			quit <- syscall.SIGINT
		}
	}()

	// Блокируем до получения сигнала
	<-quit
	s.logger.Info("Shutting down server...")

	// Создаем контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Warn("Server forced to shutdown", "why", err)
		return err
	}

	s.logger.Info("Server gracefully stopped")
	return nil
}

// --------- Helpers --------- \\

// WriteJSON отправляет данные в формате JSON с указанным HTTP статусом.
// Автоматически устанавливает правильный Content-Type заголовок.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// apiFunc определяет сигнатуру функций-обработчиков API,
// которые возвращают ошибку для централизованной обработки.
type apiFunc func(http.ResponseWriter, *http.Request) error

// APIError представляет структуру ошибки для ответов API.
type APIError struct {
	Error string `json:"error"`
}

// makeHTTPHandleFunc преобразует apiFunc в стандартный http.HandlerFunc,
// добавляя унифицированную обработку ошибок.
func (s *APIServer) makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Логирование ошибки
			s.logger.Info("Error handling request", "error", err)

			// По умолчанию используем BadRequest, но в будущем здесь можно
			// добавить логику определения правильного кода статуса на основе типа ошибки
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
