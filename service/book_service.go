package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type BookService struct {
	db *sql.DB
}

// NewBookService creates a new PostgreSQL book service
func NewBookService(db *sql.DB) *BookService {
	return &BookService{
		db: db,
	}
}

// GetAllBook retrieves all book items from PostgreSQL
func (s *BookService) GetAllBook() ([]model.Book, error) {
	query := `
		SELECT id, nama_barang, stok, terjual, harga, harga_pasar, tokped_keyword, keterangan
		FROM books 
		ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query book: %v", err)
	}
	defer rows.Close()

	var inventories []model.Book
	for rows.Next() {
		var inv model.Book
		err := rows.Scan(
			&inv.ID,
			&inv.NamaBarang,
			&inv.Stok,
			&inv.Terjual,
			&inv.Harga,
			&inv.HargaPasar,
			&inv.TokpedKeyword,
			&inv.Keterangan,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book row: %v", err)
		}
		inventories = append(inventories, inv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating book rows: %v", err)
	}

	return inventories, nil
}

// AddBook adds a new book item to PostgreSQL
func (s *BookService) AddBook(inv model.Book) error {
	query := `
		INSERT INTO books (nama_barang, stok, harga, harga_pasar, tokped_keyword, keterangan)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query,
		inv.NamaBarang,
		inv.Stok,
		inv.Harga,
		inv.HargaPasar,
		inv.TokpedKeyword,
		inv.Keterangan,
	)

	if err != nil {
		return fmt.Errorf("failed to insert book: %v", err)
	}

	log.Printf("Added book item: %s", inv.NamaBarang)
	return nil
}

// UpdateBook updates an existing book item in PostgreSQL
func (s *BookService) UpdateBook(id int, inv model.Book) error {
	query := `
		UPDATE books 
		SET nama_barang = $1, stok = $2, terjual = $3,
		    harga = $4, harga_pasar = $5, tokped_keyword = $6, 
		    keterangan = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8`

	result, err := s.db.Exec(query,
		inv.NamaBarang,
		inv.Stok,
		inv.Terjual,
		inv.Harga,
		inv.HargaPasar,
		inv.TokpedKeyword,
		inv.Keterangan,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update book: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}

	log.Printf("Updated book item ID %d: %s", id, inv.NamaBarang)
	return nil
}

// DeleteBook deletes an book item from PostgreSQL
func (s *BookService) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}

	log.Printf("Deleted book item ID %d", id)
	return nil
}

// InitializeSheet creates headers in the database (not needed for PostgreSQL)
func (s *BookService) InitializeSheet() error {
	// Headers are already created in the database schema
	log.Println("Database schema is ready")
	return nil
}

// UpdateMarketPrice updates the market price for a specific item
func (s *BookService) UpdateMarketPrice(id int, price int) error {
	query := `UPDATE books SET harga_pasar = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := s.db.Exec(query, price, id)
	if err != nil {
		return fmt.Errorf("failed to update market price: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}

	return nil
}

// SimulatePriceScraping simulates price scraping from Tokopedia
func SimulatePriceScraping(keyword string, basePrice int) int {
	if keyword == "" {
		return basePrice
	}

	// Simulate AI validation and price normalization
	// In a real implementation, this would call Tokopedia API and Gemini AI

	// For simulation, we'll use a simple algorithm
	// Generate a price between 80% to 120% of base price
	variation := 0.8 + float64(len(keyword)%40)/100.0
	scrapedPrice := int(float64(basePrice) * variation)

	// Simulate AI validation - ensure price is reasonable
	if scrapedPrice < basePrice*50/100 {
		scrapedPrice = basePrice * 85 / 100
	} else if scrapedPrice > basePrice*150/100 {
		scrapedPrice = basePrice * 115 / 100
	}

	return scrapedPrice
}
