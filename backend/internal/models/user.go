package models

import (
	"time"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	SessionId             string     `json:"session_id"`
	AccessToken           string     `json:"access_token"`
	RefreshToken          string     `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time  `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time  `json:"refresh_token_expires_at"`
	User                  GetUserRes `json:"user"`
}

type RenewAccessTokenReq struct {
	RefershToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type GetUserRes struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type CreateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UpdateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(r *CreateUserReq) *User {
	return &User{
		Name:      r.Name,
		Email:     r.Email,
		IsAdmin:   r.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) UpdateFromReq(req *UpdateUserReq) {
	u.Name = req.Name
	u.Email = req.Email
	u.Password = req.Password
	u.IsAdmin = req.IsAdmin
	u.UpdatedAt = time.Now()
}
