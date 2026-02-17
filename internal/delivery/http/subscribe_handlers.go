package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
	"github.com/opusdvs/DonWeather-ms-subscribe/internal/usecase"
)

type SubscribeHandlers struct {
	subscribeService usecase.SubscribeService
	ctx              context.Context
}

func NewSubscribeHandlers(subscribeService usecase.SubscribeService, ctx context.Context) SubscribeDelivery {
	return &SubscribeHandlers{subscribeService: subscribeService, ctx: ctx}
}

func (h *SubscribeHandlers) CreatePendingSubscribe(w http.ResponseWriter, r *http.Request) {
	var pendingSubscribe domain.PendingSubscribe
	if err := json.NewDecoder(r.Body).Decode(&pendingSubscribe); err != nil {
		log.Printf("Failed to decode pending subscribe: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := h.subscribeService.CreatePendingSubscribe(ctx, pendingSubscribe); err != nil {
		log.Printf("Failed to create pending subscribe: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Pending subscribe created"))
}

func (h *SubscribeHandlers) CreateSubscribe(w http.ResponseWriter, r *http.Request) {
	var subscribe domain.Subscribe
	if err := json.NewDecoder(r.Body).Decode(&subscribe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Failed to decode subscribe: %v", err)
		return
	}
	ctx := r.Context()
	if err := h.subscribeService.CreateSubscribe(ctx, subscribe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create subscribe: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscribe created"))
}
