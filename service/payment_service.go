package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{db: db}
}

// CreatePaymentDraft membuat draft pembayaran dari cart user
func (s *PaymentService) CreatePaymentDraft(userID int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Ambil semua cart user
	rows, err := tx.Query(`
		SELECT id, book_id, jumlah, harga
		FROM cart
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil cart user: %v", err)
	}
	defer rows.Close()

	var totalHarga int
	var totalItem int
	var carts []model.Cart
	for rows.Next() {
		var c model.Cart
		if err := rows.Scan(&c.ID, &c.BookID, &c.Jumlah, &c.Harga); err != nil {
			return 0, fmt.Errorf("gagal membaca cart: %v", err)
		}
		totalHarga += c.Harga
		totalItem += c.Jumlah
		carts = append(carts, c)
	}

	if len(carts) == 0 {
		return 0, fmt.Errorf("cart user kosong, tidak bisa membuat pembayaran")
	}

	// Insert ke tabel payment (status = 0 draft)
	queryPayment := `
		INSERT INTO payment (user_id, total_item, total_harga, status, created_at)
		VALUES ($1, $2, $3, 0, $4)
		RETURNING id
	`
	var paymentID int
	err = tx.QueryRow(queryPayment, userID, totalItem, totalHarga, time.Now()).Scan(&paymentID)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat draft payment: %v", err)
	}

	// Insert ke payment_detail
	for _, c := range carts {
		// hati-hati pembagian: harga_satuan = total_harga_cart / jumlah (jika jumlah>0)
		hargaSatuan := 0
		if c.Jumlah > 0 {
			hargaSatuan = c.Harga / c.Jumlah
		}
		totalHargaItem := c.Harga

		queryDetail := `
			INSERT INTO payment_detail (payment_id, book_id, jumlah, harga_satuan, total_harga, status, created_at)
			VALUES ($1, $2, $3, $4, $5, 0, $6)
		`
		_, err := tx.Exec(queryDetail, paymentID, c.BookID, c.Jumlah, hargaSatuan, totalHargaItem, time.Now())
		if err != nil {
			return 0, fmt.Errorf("gagal membuat detail payment: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("gagal commit payment draft: %v", err)
	}

	log.Printf("Draft payment #%d dibuat untuk user %d total_item=%d total_harga=%d", paymentID, userID, totalItem, totalHarga)
	return paymentID, nil
}

// ConfirmPayment menandai pembayaran sebagai sukses
func (s *PaymentService) ConfirmPayment(paymentID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", err)
	}
	defer tx.Rollback()

	// Update status payment
	_, err = tx.Exec(`
		UPDATE payment
		SET status = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, paymentID)
	if err != nil {
		return fmt.Errorf("gagal mengupdate payment: %v", err)
	}

	// Update detail
	_, err = tx.Exec(`
		UPDATE payment_detail
		SET status = 1
		WHERE payment_id = $1
	`, paymentID)
	if err != nil {
		return fmt.Errorf("gagal mengupdate payment detail: %v", err)
	}

	// Ambil user_id untuk hapus cart
	var userID int
	err = tx.QueryRow(`SELECT user_id FROM payment WHERE id = $1`, paymentID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("gagal mengambil user_id dari payment: %v", err)
	}

	// Hapus cart user
	_, err = tx.Exec(`DELETE FROM cart WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("gagal menghapus cart user: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit konfirmasi payment: %v", err)
	}

	log.Printf("Payment #%d dikonfirmasi dan cart user %d dikosongkan", paymentID, userID)
	return nil
}

// CancelPayment menandai pembayaran sebagai batal
func (s *PaymentService) CancelPayment(paymentID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi cancel: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE payment
		SET status = 2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, paymentID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan payment: %v", err)
	}

	_, err = tx.Exec(`
		UPDATE payment_detail
		SET status = 2
		WHERE payment_id = $1
	`, paymentID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan detail payment: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit cancel payment: %v", err)
	}

	log.Printf("Payment #%d dibatalkan", paymentID)
	return nil
}

// GetPaymentsByUser mengambil semua riwayat pembayaran user
func (s *PaymentService) GetPaymentsByUser(userID int) ([]model.Payment, error) {
	query := `
		SELECT id, user_id, total_item, total_harga, status, created_at, updated_at
		FROM payment
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat payment: %v", err)
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.TotalItem, &p.TotalHarga, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca data payment: %v", err)
		}
		payments = append(payments, p)
	}

	return payments, nil
}
