package db

import "github.com/rx3lixir/agg-api/internal/models"

type Storage interface {
	CreateAccount(*models.Account) error
	DeleteAccount(int) error
	UpdateAccount(*models.Account) error
	GetAccountByID(int) (*models.Account, error)
}
