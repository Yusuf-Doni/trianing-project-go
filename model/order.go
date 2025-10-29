package model

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID        int             `json:"id" db:"id"`
	UserID    int             `json:"user_id" db:"user_id"`
	BookID    int             `json:"book_id" db:"book_id"`
	Item      json.RawMessage `json:"item" db:"item"`
	Harga     int             `json:"harga" db:"harga"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
