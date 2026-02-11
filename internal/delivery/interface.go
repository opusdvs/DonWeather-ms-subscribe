package delivery

import "net/http"

type SubscribeDelivery interface {
	CreateSubscribe(w http.ResponseWriter, r *http.Request)
	GetAllSubscribes(w http.ResponseWriter, r *http.Request)
	GetSubscribeById(w http.ResponseWriter, r *http.Request)
	UpdateSubscribe(w http.ResponseWriter, r *http.Request)
	DeleteSubscribe(w http.ResponseWriter, r *http.Request)
}
