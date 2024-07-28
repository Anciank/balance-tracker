package models

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        int          `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Token     string       `json:"token"`
}
