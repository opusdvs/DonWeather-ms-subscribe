package domain

type Subscribe struct {
	ID         string  `json:"id"`
	TelegramID int64   `json:"telegram_id"`
	Token      string  `json:"token"`
	City       string  `json:"city"`
	Filters    Filters `json:"filters"`
}

type PendingSubscribe struct {
	ID         string  `json:"id"`
	TelegramID int64   `json:"telegram_id"`
	Token      string  `json:"token"`
	City       string  `json:"city"`
	Filters    Filters `json:"filters"`
}

type Filters struct {
	Humidity    bool `json:"humidity"`
	Temperature bool `json:"temperature"`
	WindSpeed   bool `json:"wind_speed"`
	Pressure    bool `json:"pressure"`
}
