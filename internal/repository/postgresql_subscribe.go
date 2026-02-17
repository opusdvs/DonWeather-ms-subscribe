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

func (r *PostgresqlSubscribeRepository) Save(ctx context.Context, sub domain.Subscribe) error {
	query := `
		INSERT INTO subscribe (token, telegram_id, city, filters)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, sub.Token, sub.TelegramID, sub.City, sub.Filters).Scan(&sub.ID)
}
