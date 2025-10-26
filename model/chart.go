package model

type Chart struct {
	ID     int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
	BookID int `json:"book_id" db:"book_id"`
	Jumlah int `json:"jumlah" db:"jumlah"`
	Harga  int `json:"harga" db:"harga"`
}
