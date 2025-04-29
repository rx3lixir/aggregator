package models

import "time"

// Category представляет категорию событий
type Category struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCategoryReq представляет запрос на создание новой категории
type CreateCategoryReq struct {
	Name string `json:"name"`
}

// NewCategory создает новую категорию из запроса
func NewCategory(req *CreateCategoryReq) *Category {
	return &Category{
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
