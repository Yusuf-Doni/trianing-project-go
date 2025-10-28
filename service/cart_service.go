package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type CartService struct {
	db *sql.DB
}

func NewCartService(db *sql.DB) *CartService {
	return &CartService{db: db}
}

// AddToCart menambahkan produk ke keranjang user
// Jika produk sudah ada, maka jumlah dan harga akan diupdate
func (s *CartService) AddToCart(cart model.Cart) error {
	// Cek apakah user sudah punya produk ini di cart
	queryCheck := `
		SELECT id, jumlah, harga
		FROM cart
		WHERE user_id = $1 AND book_id = $2
	`
	var existingID, existingJumlah, existingHarga int
	err := s.db.QueryRow(queryCheck, cart.UserID, cart.BookID).Scan(&existingID, &existingJumlah, &existingHarga)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("gagal mengecek cart existing: %v", err)
	}

	// Jika produk belum ada → insert baru
	if err == sql.ErrNoRows {
		queryInsert := `
			INSERT INTO cart (user_id, book_id, jumlah, harga)
			VALUES ($1, $2, $3, $4)
		`
		_, err := s.db.Exec(queryInsert, cart.UserID, cart.BookID, cart.Jumlah, cart.Harga)
		if err != nil {
			return fmt.Errorf("gagal menambahkan ke cart: %v", err)
		}
		log.Printf("Produk ID %d ditambahkan ke cart user %d (jumlah: %d)", cart.BookID, cart.UserID, cart.Jumlah)
		return nil
	}

	// Jika produk sudah ada → update jumlah dan harga total
	newJumlah := existingJumlah + cart.Jumlah
	newHarga := existingHarga + cart.Harga

	queryUpdate := `
		UPDATE cart
		SET jumlah = $1, harga = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err = s.db.Exec(queryUpdate, newJumlah, newHarga, existingID)
	if err != nil {
		return fmt.Errorf("gagal mengupdate jumlah cart: %v", err)
	}

	log.Printf("Produk ID %d di cart user %d diupdate: jumlah %d → %d", cart.BookID, cart.UserID, existingJumlah, newJumlah)
	return nil
}

// GetCartByUser mengambil semua item cart milik user tertentu
func (s *CartService) GetCartByUser(userID int) ([]model.Cart, error) {
	query := `
		SELECT id, user_id, book_id, jumlah, harga, created_at, updated_at
		FROM cart
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil cart user: %v", err)
	}
	defer rows.Close()

	var carts []model.Cart
	for rows.Next() {
		var c model.Cart
		err := rows.Scan(&c.ID, &c.UserID, &c.BookID, &c.Jumlah, &c.Harga, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data cart: %v", err)
		}
		carts = append(carts, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi cart: %v", err)
	}

	return carts, nil
}

// RemoveFromCart menghapus item dari cart berdasarkan ID cart
func (s *CartService) RemoveFromCart(cartID int) error {
	query := `DELETE FROM cart WHERE id = $1`
	result, err := s.db.Exec(query, cartID)
	if err != nil {
		return fmt.Errorf("gagal menghapus item cart: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa hasil penghapusan: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item cart dengan ID %d tidak ditemukan", cartID)
	}

	log.Printf("Item cart ID %d berhasil dihapus", cartID)
	return nil
}

// UpdateCartQuantity mengubah jumlah produk di cart
func (s *CartService) UpdateCartQuantity(cartID int, newJumlah int) error {
	var currentJumlah, hargaSatuan int

	query := `
		SELECT jumlah, (harga / jumlah) AS harga_satuan
		FROM cart
		WHERE id = $1
	`
	err := s.db.QueryRow(query, cartID).Scan(&currentJumlah, &hargaSatuan)
	if err != nil {
		return fmt.Errorf("cart tidak ditemukan: %v", err)
	}

	if newJumlah <= 0 {
		// Hapus item kalau jumlah jadi 0
		_, err := s.db.Exec(`DELETE FROM cart WHERE id = $1`, cartID)
		if err != nil {
			return fmt.Errorf("gagal menghapus cart: %v", err)
		}
		return nil
	}

	newHarga := hargaSatuan * newJumlah

	_, err = s.db.Exec(`
		UPDATE cart
		SET jumlah = $1, harga = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`, newJumlah, newHarga, cartID)

	if err != nil {
		return fmt.Errorf("gagal memperbarui jumlah cart: %v", err)
	}

	return nil
}
