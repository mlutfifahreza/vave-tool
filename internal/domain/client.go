package domain

import (
	"context"
	"time"
)

type Client struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ClientRepository interface {
	GetByUsername(ctx context.Context, username string) (*Client, error)
	GetByID(ctx context.Context, id string) (*Client, error)
}
