package entity

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	Age  *int   `json:"age" validate:"required"`
}
