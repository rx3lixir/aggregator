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
)

// APIServer представляет сервер API с настройками и обработчиками.
type APIServer struct {
	listenAddr string
	logger     logger.Logger
	server     *http.Server
	store      db.Storage
}

// NewAPIServer создает новый экземпляр APIServer с указанным адресом.
func NewAPIServer(listenAddr string, log logger.Logger, store db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		logger:     log,
		store:      store,
	}
}

// Run запускает HTTP сервер с настроенными маршрутами.
// Также настраивает корректное завершение работы при получении сигнала завершения.
func (s *APIServer) Run() error {
	router := chi.NewRouter()

	// Настройка маршрутов для работы с аккаунтами
	router.Route("/account", func(r chi.Router) {
		r.Get("/", s.makeHTTPHandleFunc(s.handleGetAccount))
		r.Get("/{id}", s.makeHTTPHandleFunc(s.handleGetAccountByID))
		r.Post("/", s.makeHTTPHandleFunc(s.handleCreateAccount))
		r.Delete("/{id}", s.makeHTTPHandleFunc(s.handleDeleteAccount))
	})

	// Создаем HTTP сервер с настроенными параметрами
	s.server = &http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}

	// Запускаем сервер в отдельной горутине
	serverErrors := make(chan error, 1)
	go func() {
		s.logger.Info("API server starting", "address", s.listenAddr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		s.logger.Info("Shutting down server due to signal...")
	case err := <-serverErrors:
		s.logger.Error("Server error occured", "error", err)
		return err
	}

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
