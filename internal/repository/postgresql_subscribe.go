package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
)

type PostgresqlSubscribeRepository struct {
	db *sql.DB
}

func NewPostgresqlSubscribeRepository(db *sql.DB) *PostgresqlSubscribeRepository {
	return &PostgresqlSubscribeRepository{db: db}
}

func (r *PostgresqlSubscribeRepository) Create(ctx context.Context, subscribe domain.Subscribe) (string, error) {
	filters, err := json.Marshal(subscribe.Filters)
	if err != nil {
		return "", err
	}
	query := `
		INSERT INTO subscribe (telegram_id, city, filters)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id string
	err = r.db.QueryRowContext(ctx, query, subscribe.TelegramID, subscribe.City, filters).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *PostgresqlSubscribeRepository) GetAll(ctx context.Context) ([]domain.Subscribe, error) {
	query := `
		SELECT id, telegram_id, city, filters
		FROM subscribe
	`
	var subscribes []domain.Subscribe
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var subscribe domain.Subscribe
		err = rows.Scan(&subscribe.ID, &subscribe.TelegramID, &subscribe.City, &subscribe.Filters)
		if err != nil {
			return nil, err
		}
		subscribes = append(subscribes, subscribe)
	}
	return subscribes, nil
}

func (r *PostgresqlSubscribeRepository) GetById(ctx context.Context, id string) (domain.Subscribe, error) {
	query := `
		SELECT id, telegram_id, city, filters
		FROM subscribe
		WHERE id = $1
	`
	var subscribe domain.Subscribe
	err := r.db.QueryRowContext(ctx, query, id).Scan(&subscribe.ID, &subscribe.TelegramID, &subscribe.City, &subscribe.Filters)
	if err != nil {
		return domain.Subscribe{}, err
	}
	return subscribe, nil
}

func (r *PostgresqlSubscribeRepository) Update(ctx context.Context, subscribe domain.Subscribe) (string, error) {
	query := `
		UPDATE subscribe
		SET telegram_id = $1, city = $2, filters = $3
		WHERE id = $4
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query, subscribe.TelegramID, subscribe.City, subscribe.Filters, subscribe.ID).Scan(&id)
	if err != nil {
		return "", err
	}
	return subscribe.ID, nil
}

func (r *PostgresqlSubscribeRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM subscribe
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
