package entity

type Sign struct {
	Key  string `json:"key" validate:"required"`
	Text string `json:"text" validate:"required"`
}
