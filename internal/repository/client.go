package repository

import (
	"context"
	"database/sql"

	"github.com/vave-tool/internal/domain"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetByUsername(ctx context.Context, username string) (*domain.Client, error) {
	query := `
		SELECT id, name, username, password, is_active, created_at, updated_at
		FROM clients
		WHERE username = $1 AND is_active = true
	`

	var client domain.Client
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&client.ID,
		&client.Name,
		&client.Username,
		&client.Password,
		&client.IsActive,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	query := `
		SELECT id, name, username, password, is_active, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	var client domain.Client
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&client.ID,
		&client.Name,
		&client.Username,
		&client.Password,
		&client.IsActive,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &client, nil
}
