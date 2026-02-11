package repository

import (
	"context"
	"database/sql"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
)

type PostgresqlSubscribeRepository struct {
	db *sql.DB
}

func NewPostgresqlSubscribeRepository(db *sql.DB) *PostgresqlSubscribeRepository {
	return &PostgresqlSubscribeRepository{db: db}
}

func (r *PostgresqlSubscribeRepository) Create(ctx context.Context, subscribe domain.Subscribe) (string, error) {
	return "", nil
}

func (r *PostgresqlSubscribeRepository) GetAll(ctx context.Context) ([]domain.Subscribe, error) {
	return nil, nil
}

func (r *PostgresqlSubscribeRepository) GetById(ctx context.Context, id string) (domain.Subscribe, error) {
	return domain.Subscribe{}, nil
}

func (r *PostgresqlSubscribeRepository) Update(ctx context.Context, subscribe domain.Subscribe) error {
	return nil
}

func (r *PostgresqlSubscribeRepository) Delete(ctx context.Context, id string) error {
	return nil
}
