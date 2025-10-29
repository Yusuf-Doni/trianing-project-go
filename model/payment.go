package model

import "time"

type Payment struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalItem  int       `json:"total_item"`
	TotalHarga int       `json:"total_harga"`
	Status     int       `json:"status"` // 0=draft, 1=sukses, 2=batal
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
