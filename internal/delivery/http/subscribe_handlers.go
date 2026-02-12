package delivery

import (
	"encoding/json"
	"errors"
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
	var subscribe domain.Subscribe
	err := json.NewDecoder(r.Body).Decode(&subscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ValidateSubscribe(subscribe); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.subscribeService.CreateSubscribe(r.Context(), subscribe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
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

func (h *SubscribeHandlers) SetTelegramID(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	telegramID := r.URL.Query().Get("telegram_id")
	h.subscribeService.SetTelegramID(r.Context(), token, telegramID)
}

func ValidateSubscribe(subscribe domain.Subscribe) error {
	if subscribe.Token == "" {
		return errors.New("token is required")
	}
	if subscribe.City == "" {
		return errors.New("city is required")
	}
	if subscribe.Filters == (domain.Filters{}) {
		return errors.New("filters are required")
	}
	return nil
}
