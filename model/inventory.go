package model

import "time"

// Inventory represents the inventory structure matching Google Sheets
type Inventory struct {
	ID            int       `json:"id" db:"id"`
	NamaBarang    string    `json:"nama_barang" db:"nama_barang"`
	StokDimiliki  int       `json:"stok_dimiliki" db:"stok_dimiliki"`
	StokTerjual   int       `json:"stok_terjual" db:"stok_terjual"`
	StokMasuk     int       `json:"stok_masuk" db:"stok_masuk"`
	HargaJual     int       `json:"harga_jual" db:"harga_jual"`
	HargaBeli     int       `json:"harga_beli" db:"harga_beli"`
	HargaPasar    int       `json:"harga_pasar" db:"harga_pasar"`
	TokpedKeyword string    `json:"tokped_keyword" db:"tokped_keyword"`
	Keterangan    string    `json:"keterangan" db:"keterangan"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
