package model

import "time"

// Cart merepresentasikan item di keranjang belanja
type Cart struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	BookID    int       `json:"book_id" db:"book_id"`
	Jumlah    int       `json:"jumlah" db:"jumlah"`
	Harga     int       `json:"harga" db:"harga"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
