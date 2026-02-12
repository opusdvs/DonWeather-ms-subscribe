package domain

import "context"

type SubscribeRepository interface {
	Create(ctx context.Context, subscribe Subscribe) (string, error)
	GetAll(ctx context.Context) ([]Subscribe, error)
	GetById(ctx context.Context, id string) (Subscribe, error)
	Update(ctx context.Context, subscribe Subscribe) (string, error)
	Delete(ctx context.Context, id string) error
	SetTelegramID(ctx context.Context, id string, telegramID string) (string, error)
}
