package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rx3lixir/agg-api/internal/lib/password"
	"github.com/rx3lixir/agg-api/internal/models"
)

// handleGetAccount обрабатывает GET запросы на /account.
// Возвращает информацию обо всех аккаунтах.
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// Этот хэндлер нужно помочь реализовать
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// Получаем ID из параметра URL
	idString := chi.URLParam(r, "id")

	// Преобразуем строку ID в число
	var id int

	if _, err := fmt.Sscanf(idString, "%d", &id); err != nil {
		return fmt.Errorf("invalid account ID: %s", idString)
	}

	// Достаем аккаунт из хранилища
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account with ID: %d - %w", id, err)
	}

	// Если аккаунт не найден - возвращаем 404
	if account == nil {
		return WriteJSON(w, http.StatusNotFound, APIError{Error: fmt.Sprintf("account with ID %d not found", id)})
	}

	// Отправляем на клиент данные
	return WriteJSON(w, http.StatusOK, account)
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
	return fmt.Errorf("not implemented yet")
}
