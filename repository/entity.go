package repository

import "time"

type OrganizationEntity struct {
	Id        int        `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Level     int        `json:"level" db:"level"`
	ParentId  *int       `json:"parent_id" db:"parent_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
