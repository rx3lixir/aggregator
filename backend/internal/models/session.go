package models

import "time"

type Session struct {
	Id           string
	UserEmail    string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
