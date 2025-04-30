package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rx3lixir/agg-api/internal/models"
)

// handleGetEvents возвращает информацию обо всех событиях
func (s *APIServer) handleGetEvents(w http.ResponseWriter, r *http.Request) error {
	events, err := s.store.GetEvents(s.dbContext)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, events)
}

// handleGetEventByID возвращает событие с переданным id
func (s *APIServer) handleGetEventById(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDFromURL(r, "id")
	if err != nil {
		return err
	}

	event, err := s.store.GetEventByID(s.dbContext, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return WriteJSON(w, http.StatusNotFound, APIError{Error: err.Error()})
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, event)
}

// handleCreateEvent создает ивент в бд
func (s *APIServer) handleCreateEvent(w http.ResponseWriter, r *http.Request) error {
	createEventReq := new(models.CreateEventReq)

	if err := json.NewDecoder(r.Body).Decode(createEventReq); err != nil {
		return err
	}

	// Базовая валидация
	if createEventReq.Name == "" {
		return fmt.Errorf("event name is required")
	}

	newEvent := models.NewEvent(createEventReq)
	if err := s.store.CreateEvent(s.dbContext, newEvent); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, newEvent)
}

// handleUpdateEvent обновляет переданное событие
func (s *APIServer) handleUpdateEvent(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDFromURL(r, "id")
	if err != nil {
		return err
	}

	// Получаем текущее событие
	event, err := s.store.GetEventByID(s.dbContext, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return WriteJSON(w, http.StatusNotFound, APIError{Error: err.Error()})
		}
		return err
	}

	// Декодируем данные обновления
	updateReq := new(models.UpdateEventReq)
	if err := json.NewDecoder(r.Body).Decode(updateReq); err != nil {
		return err
	}

	// Применение обновлений
	event.UpdateFromReq(updateReq)

	// Сохраняем изменения
	if err := s.store.UpdateEvent(s.dbContext, event); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, event)
}

// handleDeleteEvent удаляет указанное событие
func (s *APIServer) handleDeleteEvent(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDFromURL(r, "id")
	if err != nil {
		return err
	}

	if err := s.store.DeleteEvent(s.dbContext, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return WriteJSON(w, http.StatusNotFound, APIError{Error: err.Error()})
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("event %d successfully deleted", id),
	})
}

// handleGetEventsByCategory возвращает список событий для указанной категории
func (s *APIServer) handleGetEventsByCategory(w http.ResponseWriter, r *http.Request) error {
	categoryId, err := parseIDFromURL(r, "categoryId")
	if err != nil {
		return err
	}

	events, err := s.store.GetEventsByCategory(s.dbContext, categoryId)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, events)
}

func parseIDFromURL(r *http.Request, paramName string) (int, error) {
	idParam := chi.URLParam(r, paramName)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: %v", paramName, idParam)
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid %s: must be positive, got %d", paramName, id)
	}

	return id, nil
}
