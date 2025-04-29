package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rx3lixir/agg-api/internal/lib/password"
	"github.com/rx3lixir/agg-api/internal/models"
)

// handleGetUsers возвращает информацию обо всех аккаунтах.
func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetUsers(s.dbContext)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// handleGetAccount возвращает пользователя с переданным id.
func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseAndValidateID(r)
	if err != nil {
		return err
	}

	// Достаем аккаунт из хранилища
	user, err := s.store.GetUserByID(s.dbContext, id)
	if err != nil {
		return err
	}

	// Если аккаунт не найден - возвращаем 404
	if user == nil {
		return WriteJSON(w, http.StatusNotFound, APIError{Error: fmt.Sprintf("account with ID %d not found", id)})
	}

	// Отправляем на клиент данные
	return WriteJSON(w, http.StatusOK, user)
}

// handleCreateUser Cоздает новый аккаунт на основе данных из тела запроса.
func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(models.CreateUserReq)

	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	newUser := models.NewUser(createUserReq)

	hash, err := password.Hash(createUserReq.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	newUser.Password = hash

	if err := s.store.CreateUser(s.dbContext, newUser); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, newUser)
}

// handleDeleteUser Удаляет указанный аккаунт.
func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := parseAndValidateID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{
			Error: err.Error(),
		})
	}

	if err := s.store.DeleteUser(s.dbContext, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return WriteJSON(w, http.StatusNotFound, APIError{
				Error: err.Error(),
			})
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("user %d successfully deleted", id)})
}

// ParseAndValidateID извлекает ID из параметров URL и валидирует его
func parseAndValidateID(r *http.Request) (int, error) {
	idString := chi.URLParam(r, "id")

	var id int

	if _, err := fmt.Sscanf(idString, "%d", &id); err != nil {
		return 0, fmt.Errorf("invalid account ID format: %s", idString)
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid account ID: must be a postive, got %d", id)
	}

	return id, nil
}
