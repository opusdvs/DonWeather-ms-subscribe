package usecase

import (
	"context"
	"log"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
)

type SubscribeService struct {
	subscribeRepository        domain.SubscribeRepository
	pendingSubscribeRepository domain.PendengSubscribeRepository
}

func NewSubscribeService(
	subscribeRepository domain.SubscribeRepository,
	pendingSubscribeRepository domain.PendengSubscribeRepository,
) *SubscribeService {
	return &SubscribeService{
		subscribeRepository:        subscribeRepository,
		pendingSubscribeRepository: pendingSubscribeRepository,
	}
}

func (s *SubscribeService) CreatePendingSubscribe(ctx context.Context, pendingSub domain.PendingSubscribe) error {
	return s.pendingSubscribeRepository.Save(ctx, pendingSub)
}

func (s *SubscribeService) GetPendingSubscribe(ctx context.Context, token string) (domain.PendingSubscribe, error) {
	return s.pendingSubscribeRepository.Get(ctx, token)
}

func (s *SubscribeService) DeletePendingSubscribe(ctx context.Context, token string) error {
	return s.pendingSubscribeRepository.Delete(ctx, token)
}

func (s *SubscribeService) CreateSubscribe(ctx context.Context, subscribe domain.Subscribe) error {
	pendingSubscribe, err := s.pendingSubscribeRepository.Get(ctx, subscribe.Token)
	if err != nil {

		return err
	}
	finalSubscribe := domain.Subscribe{
		ID:         pendingSubscribe.ID,
		TelegramID: subscribe.TelegramID,
		City:       pendingSubscribe.City,
		Filters:    pendingSubscribe.Filters,
		Token:      pendingSubscribe.Token,
	}
	log.Printf("Final subscribe: %+v", finalSubscribe)
	if err := s.subscribeRepository.Save(ctx, finalSubscribe); err != nil {
		log.Printf("Failed to save subscribe: %v", err)
		return err
	}
	log.Printf("Subscribe saved")
	if err := s.pendingSubscribeRepository.Delete(ctx, pendingSubscribe.Token); err != nil {
		log.Printf("Failed to delete pending subscribe: %v", err)
	}
	log.Printf("Pending subscribe deleted")
	return nil
}
