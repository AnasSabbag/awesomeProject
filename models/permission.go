package models

type Permission struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" validate:"unique"`
	Description string `json:"description" db:"description" `
	IsDeleted   bool   `json:"is_deleted" db:"is_deleted"`
}

type PermissionPayload struct {
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
}

type NewRolePayload struct {
	Name         string `json:"name" db:"name" validate:"required"`
	Description  string `json:"description" db:"description"`
	PermissionID []int  `json:"permission_id" db:"permission_id"`
}
