package entity

type Value struct {
	Key   string `json:"key" validate:"required"`
	Value int    `json:"value" validate:"required"`
}
