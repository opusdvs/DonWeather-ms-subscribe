package usecase

import (
	"context"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
)

type SubscribeService struct {
	subscribeRepository domain.SubscribeRepository
}

func NewSubscribeService(subscribeRepository domain.SubscribeRepository) *SubscribeService {
	return &SubscribeService{subscribeRepository: subscribeRepository}
}

func (s *SubscribeService) CreateSubscribe(ctx context.Context, subscribe domain.Subscribe) (string, error) {
	return s.subscribeRepository.Create(ctx, subscribe)
}

func (s *SubscribeService) GetAllSubscribes(ctx context.Context) ([]domain.Subscribe, error) {
	return s.subscribeRepository.GetAll(ctx)
}

func (s *SubscribeService) GetSubscribeById(ctx context.Context, id string) (domain.Subscribe, error) {
	return s.subscribeRepository.GetById(ctx, id)
}

func (s *SubscribeService) UpdateSubscribe(ctx context.Context, subscribe domain.Subscribe) (string, error) {
	return s.subscribeRepository.Update(ctx, subscribe)
}

func (s *SubscribeService) DeleteSubscribe(ctx context.Context, id string) error {
	return s.subscribeRepository.Delete(ctx, id)
}

func (s *SubscribeService) SetTelegramID(ctx context.Context, token string, telegramID string) (string, error) {
	return s.subscribeRepository.SetTelegramID(ctx, token, telegramID)
}
