package model

import "time"

type OrderItem struct {
	ID        int
	OrderID   int
	BookID    int
	Jumlah    int
	Harga     int
	CreatedAt time.Time
}
