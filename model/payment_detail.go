package model

import "time"

type PaymentDetail struct {
	ID          int       `json:"id"`
	PaymentID   int       `json:"payment_id"`
	BookID      int       `json:"book_id"`
	Jumlah      int       `json:"jumlah"`
	HargaSatuan int       `json:"harga_satuan"`
	TotalHarga  int       `json:"total_harga"`
	Status      int       `json:"status"` // 0=draft, 1=sukses, 2=batal
	CreatedAt   time.Time `json:"created_at"`
}
