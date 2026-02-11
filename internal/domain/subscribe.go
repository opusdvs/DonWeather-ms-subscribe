package domain

type Subscribe struct {
	ID           string `json:"id"`
	TelegramName string `json:"telegram_name"`
	City         string `json:"city"`
	Params       string `json:"params"`
}
