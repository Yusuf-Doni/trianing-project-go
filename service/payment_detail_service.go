package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type PaymentDetailService struct {
	db *sql.DB
}

func NewPaymentDetailService(db *sql.DB) *PaymentDetailService {
	return &PaymentDetailService{db: db}
}

// CreatePaymentDetail menambahkan detail pembayaran (per barang)
func (s *PaymentDetailService) CreatePaymentDetail(detail model.PaymentDetail) error {
	query := `
		INSERT INTO payment_detail (payment_id, book_id, jumlah, harga_satuan, total_harga, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
	`
	_, err := s.db.Exec(query, detail.PaymentID, detail.BookID, detail.Jumlah, detail.HargaSatuan, detail.TotalHarga, detail.Status)
	if err != nil {
		return fmt.Errorf("gagal menambahkan payment detail: %v", err)
	}
	log.Printf("Detail pembayaran untuk payment_id %d ditambahkan (book_id: %d, jumlah: %d)", detail.PaymentID, detail.BookID, detail.Jumlah)
	return nil
}

// GetPaymentDetailsByPaymentID mengambil semua detail berdasarkan ID payment
func (s *PaymentDetailService) GetPaymentDetailsByPaymentID(paymentID int) ([]model.PaymentDetail, error) {
	query := `
		SELECT id, payment_id, book_id, jumlah, harga_satuan, total_harga, status, created_at
		FROM payment_detail
		WHERE payment_id = $1
		ORDER BY id
	`

	rows, err := s.db.Query(query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pembayaran: %v", err)
	}
	defer rows.Close()

	var details []model.PaymentDetail
	for rows.Next() {
		var d model.PaymentDetail
		err := rows.Scan(&d.ID, &d.PaymentID, &d.BookID, &d.Jumlah, &d.HargaSatuan, &d.TotalHarga, &d.Status, &d.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data detail pembayaran: %v", err)
		}
		details = append(details, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterasi data payment_detail: %v", err)
	}

	return details, nil
}

// DeletePaymentDetailsByPaymentID menghapus semua detail berdasarkan payment ID (misal kalau payment dibatalkan)
func (s *PaymentDetailService) DeletePaymentDetailsByPaymentID(paymentID int) error {
	query := `DELETE FROM payment_detail WHERE payment_id = $1`
	result, err := s.db.Exec(query, paymentID)
	if err != nil {
		return fmt.Errorf("gagal menghapus detail pembayaran: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Payment detail untuk payment_id %d dihapus (%d baris)", paymentID, rowsAffected)
	return nil
}
