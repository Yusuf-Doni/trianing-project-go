package model

type Chart struct {
	ID         int    `json:"id" db:"id"`
	NamaBarang string `json:"nama_barang" db:"nama_barang"`
}
