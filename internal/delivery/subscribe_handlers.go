package delivery

import (
	"net/http"

	"github.com/opusdvs/DonWeather-ms-subscribe/internal/domain"
	"github.com/opusdvs/DonWeather-ms-subscribe/internal/usecase"
)

type SubscribeHandlers struct {
	subscribeService usecase.SubscribeService
}

func NewSubscribeHandlers(subscribeService usecase.SubscribeService) SubscribeDelivery {
	return &SubscribeHandlers{subscribeService: subscribeService}
}

func (h *SubscribeHandlers) CreateSubscribe(w http.ResponseWriter, r *http.Request) {
	h.subscribeService.CreateSubscribe(r.Context(), domain.Subscribe{})
}

func (h *SubscribeHandlers) GetAllSubscribes(w http.ResponseWriter, r *http.Request) {
	h.subscribeService.GetAllSubscribes(r.Context())
}

func (h *SubscribeHandlers) GetSubscribeById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	h.subscribeService.GetSubscribeById(r.Context(), id)
}

func (h *SubscribeHandlers) UpdateSubscribe(w http.ResponseWriter, r *http.Request) {
	h.subscribeService.UpdateSubscribe(r.Context(), domain.Subscribe{})
}

func (h *SubscribeHandlers) DeleteSubscribe(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	h.subscribeService.DeleteSubscribe(r.Context(), id)
}
