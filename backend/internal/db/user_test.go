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

// --- Тесты для DeleteUser --- \\

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

// --- Тесты для GetUsers --- \\

func TestPostgresStore_GetUsers_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	now := time.Now()

	user1 := &models.User{
		Id:        1,
		Name:      "Kotichek",
		Email:     "yulia@kisik.ru",
		Password:  "nikakihklubkov123",
		IsAdmin:   false,
		CreatedAt: now.Add(-time.Hour * 24),
		UpdatedAt: now.Add(-time.Hour),
	}

	user2 := &models.User{
		Id:        2,
		Name:      "mister_KPblC",
		Email:     "ilovecheese@mail.ru",
		Password:  "mousetrap3434",
		IsAdmin:   false,
		CreatedAt: now.Add(-time.Hour * 9),
		UpdatedAt: now.Add(-time.Hour),
	}

	expectedUsers := []*models.User{user1, user2}

	// Определяем столбцы, возвращаемые запросом SELECT *
	cols := []string{"id", "name", "email", "password", "is_admin", "created_at", "updated_at"}

	rows := pgxmock.NewRows(cols).AddRow(
		user1.Id,
		user1.Name,
		user1.Email,
		user1.Password,
		user1.IsAdmin,
		user1.CreatedAt,
		user1.UpdatedAt,
	).AddRow(
		user2.Id,
		user2.Name,
		user2.Email,
		user2.Password,
		user2.IsAdmin,
		user2.CreatedAt,
		user2.UpdatedAt,
	)

	expectedSQL := regexp.QuoteMeta("SELECT * FROM users")
	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

	store := db.NewPosgresStore(mock)
	users, err := store.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_GetUsers_EmptyResult(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	// Пустой набор результатов
	cols := []string{"id", "name", "email", "password", "is_admin", "created_at", "updated_at"}
	rows := pgxmock.NewRows(cols)

	expectedSQL := regexp.QuoteMeta("SELECT * FROM users")
	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

	store := db.NewPosgresStore(mock)
	users, err := store.GetUsers(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_GetUsers_DatabaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	expectedSQL := regexp.QuoteMeta("SELECT * FROM users")
	mock.ExpectQuery(expectedSQL).WillReturnError(fmt.Errorf("database connection error"))

	store := db.NewPosgresStore(mock)
	users, err := store.GetUsers(context.Background())

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_GetUsers_ScanError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	// Создаем строки с некорректными данными (например, неверный тип)
	cols := []string{"id", "name", "email", "password", "is_admin", "created_at", "updated_at"}
	rows := pgxmock.NewRows(cols).
		AddRow(
			"not_an_integer", // Неверный тип для id (должно быть int)
			"User1",
			"user1@example.com",
			"password1",
			false,
			time.Now(),
			time.Now(),
		)

	expectedSQL := regexp.QuoteMeta("SELECT * FROM users")
	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

	store := db.NewPosgresStore(mock)
	users, err := store.GetUsers(context.Background())

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Тесты для UpdateUser --- \\

func TestPostgresStore_UpdateUser_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	userID := 5
	now := time.Now()
	updatedUser := &models.User{
		Id:       userID,
		Name:     "NewName",
		Email:    "newemail@example.com",
		Password: "newpassword",
		IsAdmin:  false,
	}

	// Ожидаем проверку существования пользователя
	existsQuery := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)")
	mock.ExpectQuery(existsQuery).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	// Ожидаем обновление пользователя
	updateQuery := `UPDATE users SET name = \$1, email = \$2, password = \$3, updated_at = NOW\(\) WHERE id = \$4 RETURNING updated_at`
	mock.ExpectQuery(updateQuery).
		WithArgs(updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Id).
		WillReturnRows(pgxmock.NewRows([]string{"updated_at"}).AddRow(now))

	store := db.NewPosgresStore(mock)
	err = store.UpdateUser(context.Background(), updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, now, updatedUser.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_UpdateUser_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	userID := 999
	updatedUser := &models.User{
		Id:       userID,
		Name:     "NewName",
		Email:    "newemail@example.com",
		Password: "newpassword",
		IsAdmin:  false,
	}

	// Ожидаем проверку существования пользователя, который не существует
	existsQuery := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)")
	mock.ExpectQuery(existsQuery).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	store := db.NewPosgresStore(mock)
	err = store.UpdateUser(context.Background(), updatedUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("user with ID %d not found", userID))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresStore_UpdateUser_DatabaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Ошибка создания мока: %v", err)
	}
	defer mock.Close()

	userID := 5
	updatedUser := &models.User{
		Id:       userID,
		Name:     "NewName",
		Email:    "newemail@example.com",
		Password: "newpassword",
		IsAdmin:  false,
	}

	// Ожидаем проверку существования пользователя
	existsQuery := regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)")
	mock.ExpectQuery(existsQuery).
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	// Ожидаем ошибку при обновлении
	updateQuery := `UPDATE users SET name = \$1, email = \$2, password = \$3, updated_at = NOW\(\) WHERE id = \$4 RETURNING updated_at`
	mock.ExpectQuery(updateQuery).
		WithArgs(updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Id).
		WillReturnError(fmt.Errorf("database error"))

	store := db.NewPosgresStore(mock)
	err = store.UpdateUser(context.Background(), updatedUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("failed to update user %d", userID))
	assert.NoError(t, mock.ExpectationsWereMet())
}
