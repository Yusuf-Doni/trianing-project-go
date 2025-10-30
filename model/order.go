package model

import (
	"time"
)

type Order struct {
	ID         int
	UserID     int
	TotalHarga int
	CreatedAt  time.Time
}
