package db_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"github.com/rx3lixir/agg-api/internal/db"
	"github.com/rx3lixir/agg-api/internal/models"
)

// --- Тест для GetUserByID --- \\

func TestPostgresStore_GetUserByID_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	expectedID := 1
	now := time.Now()
	expectedUser := &models.User{
		Id:        expectedID,
		Name:      "mr_KRbIC",
		Email:     "KRbIClovesCHEES@mail.ru",
		Password:  "1235",
		IsAdmin:   false,
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now,
	}

	cols := []string{"id", "name", "email", "password", "is_amdin", "created_at", "updated_at"}

	rows := pgxmock.NewRows(cols).AddRow(
		expectedUser.Id,
		expectedUser.Name,
		expectedUser.Email,
		expectedUser.Password,
		expectedUser.IsAdmin,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
	)

	expectedSQL := regexp.QuoteMeta("SELECT id, name, email, password, is_admin, created_at, updated_at FROM users WHERE id = $1")
	mock.ExpectQuery(expectedSQL).WithArgs(expectedID).WillReturnRows(rows)

	store := db.NewPosgresStore(mock)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user, err := store.GetUserByID(ctx, expectedID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_GetUserByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	testID := 999
	expectedSQL := regexp.QuoteMeta("SELECT id, name, email, password, is_admin, created_at, updated_at FROM users WHERE id = $1")

	// Ожидаем, что QueryRow вернет ошибку pgx.ErrNoRows
	mock.ExpectQuery(expectedSQL).
		WithArgs(testID).
		WillReturnError(pgx.ErrNoRows) // Симулируем ошибку "не найдено"

	store := db.NewPosgresStore(mock)
	user, err := store.GetUserByID(context.Background(), testID)

	// Проверяем результат
	assert.Error(t, err) // Ошибка должна быть
	// Проверяем, что ошибка содержит ожидаемое сообщение (т.к. GetUserByID форматирует ошибку)
	assert.Contains(t, err.Error(), fmt.Sprintf("user %d not found", testID))
	// Или, если бы GetUserByID возвращал обернутую ошибку: assert.ErrorIs(t, err, pgx.ErrNoRows)
	assert.Nil(t, user) // Пользователь должен быть nil

	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Тест для CreateUser --- \\

func TestPostgresStore_CreateUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	now := time.Now()
	inputUser := &models.User{
		Name:     "Govno",
		Email:    "govno@example.com",
		Password: "plainpassword",
		IsAdmin:  false,
	}
	// Ожидаемые значения, которые вернет база данных
	expectedID := 5
	expectedCreatedAt := now
	expectedUpdatedAt := now

	// Колонки, возвращаемые RETURNING
	returningCols := []string{"id", "created_at", "updated_at"}
	// Строка с возвращаемыми значениями
	returningRows := pgxmock.NewRows(returningCols).
		AddRow(expectedID, expectedCreatedAt, expectedUpdatedAt)

	// Ожидаем INSERT ... RETURNING (это QueryRow для pgx)
	// Используем `.+` для гибкого соответствия SQL
	// Важно: порядок аргументов в WithArgs должен соответствовать VALUES ($1, $2, $3, $4)
	mock.ExpectQuery(`INSERT INTO users \(name, email, password, is_admin\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, created_at, updated_at`).
		WithArgs(inputUser.Name, inputUser.Email, inputUser.Password, inputUser.IsAdmin).
		WillReturnRows(returningRows)

	store := db.NewPosgresStore(mock)
	err = store.CreateUser(context.Background(), inputUser)

	// Проверяем результат
	assert.NoError(t, err)
	// Проверяем, что поля Id, CreatedAt, UpdatedAt были установлены в inputUser
	assert.Equal(t, expectedID, inputUser.Id)
	assert.Equal(t, expectedCreatedAt, inputUser.CreatedAt)
	assert.Equal(t, expectedUpdatedAt, inputUser.UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Тест для DeleteUser --- \\

func TestPostgresStore_DeleteUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	testID := 10

	// Ожидаем DELETE (это Exec для pgx)
	expectedSQL := regexp.QuoteMeta("DELETE FROM users WHERE id = $1")
	mock.ExpectExec(expectedSQL).
		WithArgs(testID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1)) // Симулируем, что 1 строка удалена

	store := db.NewPosgresStore(mock)
	err = store.DeleteUser(context.Background(), testID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_DeleteUser_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	testID := 99

	// Ожидаем DELETE, но он не затронет ни одной строки
	expectedSQL := regexp.QuoteMeta("DELETE FROM users WHERE id = $1")
	mock.ExpectExec(expectedSQL).
		WithArgs(testID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0)) // Симулируем, что 0 строк удалено

	store := db.NewPosgresStore(mock)
	err = store.DeleteUser(context.Background(), testID)

	// Проверяем результат
	assert.Error(t, err) // Должна быть ошибка, т.к. RowsAffected == 0
	assert.Contains(t, err.Error(), fmt.Sprintf("user with ID %d not found for deletion", testID))
	assert.NoError(t, mock.ExpectationsWereMet()) // Ожидание Exec было выполнено, хоть и с 0 результатом
}
