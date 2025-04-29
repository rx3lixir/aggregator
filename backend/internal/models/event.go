package models

import (
	"time"
)

// CreateEventReq представляет запрос на создание нового события
type CreateEventReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryId  int     `json:"category_id"`
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	Image       string  `json:"image"`
	Source      string  `json:"source"`
}

// UpdateEventReq представляет запрос на обновление существующего события
type UpdateEventReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryId  int     `json:"category_id"`
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	Rating      float64 `json:"rating"`
	Image       string  `json:"image"`
	Source      string  `json:"source"`
}

// Event представляет событие в системе
type Event struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryId  int       `json:"category_id"`
	Date        string    `json:"date"`
	Time        string    `json:"time"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
	Rating      float64   `json:"rating"`
	Image       string    `json:"image"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewEvent создает новый экземпляр события на основе запроса
func NewEvent(req *CreateEventReq) *Event {
	return &Event{
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryId,
		Date:        req.Date,
		Time:        req.Time,
		Location:    req.Location,
		Price:       req.Price,
		Rating:      req.Rating,
		Image:       req.Image,
		Source:      req.Source,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateFromReq обновляет существующее событие из запроса
func (e *Event) UpdateFromReq(req *UpdateEventReq) {
	e.Name = req.Name
	e.Description = req.Description
	e.CategoryId = req.CategoryId
	e.Date = req.Date
	e.Time = req.Time
	e.Location = req.Location
	e.Price = req.Price
	e.Rating = req.Rating
	e.Image = req.Image
	e.Source = req.Source
	e.UpdatedAt = time.Now()
}
