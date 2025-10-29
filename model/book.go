package model

import "time"

// Book represents the book structure matching Google Sheets
type Book struct {
	ID            int       `json:"id" db:"id"`
	NamaBarang    string    `json:"nama_barang" db:"nama_barang"`
	Stok          int       `json:"stok" db:"stok"`
	Terjual       int       `json:"terjual" db:"terjual"`
	Harga         int       `json:"harga" db:"harga"`
	HargaPasar    int       `json:"harga_pasar" db:"harga_pasar"`
	TokpedKeyword string    `json:"tokped_keyword" db:"tokped_keyword"`
	Keterangan    string    `json:"keterangan" db:"keterangan"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
