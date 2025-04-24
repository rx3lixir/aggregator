package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rx3lixir/agg-api/internal/lib/password"
	"github.com/rx3lixir/agg-api/internal/models"
)

// handleGetAccount обрабатывает GET запросы на /account.
// Возвращает информацию о запрошенном аккаунте.
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "Hi there! You've reached the /account handler")
}

// handleCreateAccount обрабатывает POST запросы на /account.
// Создает новый аккаунт на основе данных из тела запроса.
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(models.CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := models.NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Email)

	hash, err := password.Hash(createAccountReq.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	account.PasswordHash = hash

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// handleDeleteAccount обрабатывает DELETE запросы на /account.
// Удаляет указанный аккаунт.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// TODO: Реализовать удаление аккаунта
	// 1. Получить ID аккаунта из URL или параметров
	// 2. Проверить существование аккаунта
	// 3. Удалить аккаунт
	// 4. Вернуть результат

	return fmt.Errorf("not implemented yet")
}
