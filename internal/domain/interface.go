package domain

import "context"

type PendengSubscribeRepository interface {
	Save(ctx context.Context, sub PendingSubscribe) error
	Get(ctx context.Context, token string) (PendingSubscribe, error)
	Delete(ctx context.Context, token string) error
}

type SubscribeRepository interface {
	Save(ctx context.Context, sub Subscribe) error
}
