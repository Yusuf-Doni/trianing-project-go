package model

import "time"

type Book struct {
	ID            int
	NamaBarang    string
	Stok          int
	Terjual       int
	Harga         int
	HargaPasar    int
	TokpedKeyword string
	Keterangan    string
	GambarBuku    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
