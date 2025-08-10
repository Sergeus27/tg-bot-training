package telegram

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}
