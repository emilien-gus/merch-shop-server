package models

type Purchase struct {
	ID       int    `json:"-" db:"id"`
	UserID   int    `json:"-" db:"user_id"`
	ItemName string `json:"item_name" db:"item_name"`
	Quantity int    `json:"quantity" db:"quantity"`
	Price    int    `json:"price" db:"price"`
}
