package models

import "time"

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name" db:"name" valid:"required"`
	Email             string    `json:"email" db:"email" valid:"email,required"`
	Username          string    `json:"username" db:"username" valid:"required"`
	Password          string    `json:"password" db:"password" valid:"required"`
	ConfirmedPassword string    `json:"confirmed" db:"confirmed_password" valid:"eqfield=Password"`
	RoleID            int       `json:"roleId" db:"role_id" valid:"required"`
	CreatedAt         time.Time `json:"createdAt"`
	ArchivedAt        time.Time `json:"archivedAt"`
	LastLoginAt       time.Time `json:"lastLoginAt"`
	IsAdmin           bool      `json:"isAdmin" db:"is_admin"`
	IsDeactivate      bool      `json:"isActive"`
}
type LoginCredentials struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}
