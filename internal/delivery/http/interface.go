package delivery

import "net/http"

type SubscribeDelivery interface {
	CreatePendingSubscribe(w http.ResponseWriter, r *http.Request)
	CreateSubscribe(w http.ResponseWriter, r *http.Request)
}
