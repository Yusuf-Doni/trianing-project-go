package model

import "time"

type Session struct {
	ID        string
	UserID    int
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}
