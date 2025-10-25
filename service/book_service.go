package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type InventoryService struct {
	db *sql.DB
}

// NewInventoryService creates a new PostgreSQL inventory service
func NewInventoryService(db *sql.DB) *InventoryService {
	return &InventoryService{
		db: db,
	}
}

// GetAllInventory retrieves all inventory items from PostgreSQL
func (s *InventoryService) GetAllInventory() ([]model.Inventory, error) {
	query := `
		SELECT id, nama_barang, stok_dimiliki, stok_terjual, stok_masuk, 
		       harga_jual, harga_beli, harga_pasar, tokped_keyword, keterangan
		FROM inventory 
		ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory: %v", err)
	}
	defer rows.Close()

	var inventories []model.Inventory
	for rows.Next() {
		var inv model.Inventory
		err := rows.Scan(
			&inv.ID,
			&inv.NamaBarang,
			&inv.StokDimiliki,
			&inv.StokTerjual,
			&inv.StokMasuk,
			&inv.HargaJual,
			&inv.HargaBeli,
			&inv.HargaPasar,
			&inv.TokpedKeyword,
			&inv.Keterangan,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory row: %v", err)
		}
		inventories = append(inventories, inv)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inventory rows: %v", err)
	}

	return inventories, nil
}

// AddInventory adds a new inventory item to PostgreSQL
func (s *InventoryService) AddInventory(inv model.Inventory) error {
	query := `
		INSERT INTO inventory (nama_barang, stok_dimiliki, stok_terjual, stok_masuk, 
		                     harga_jual, harga_beli, harga_pasar, tokped_keyword, keterangan)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := s.db.Exec(query,
		inv.NamaBarang,
		inv.StokDimiliki,
		inv.StokTerjual,
		inv.StokMasuk,
		inv.HargaJual,
		inv.HargaBeli,
		inv.HargaPasar,
		inv.TokpedKeyword,
		inv.Keterangan,
	)

	if err != nil {
		return fmt.Errorf("failed to insert inventory: %v", err)
	}

	log.Printf("Added inventory item: %s", inv.NamaBarang)
	return nil
}

// UpdateInventory updates an existing inventory item in PostgreSQL
func (s *InventoryService) UpdateInventory(id int, inv model.Inventory) error {
	query := `
		UPDATE inventory 
		SET nama_barang = $1, stok_dimiliki = $2, stok_terjual = $3, stok_masuk = $4,
		    harga_jual = $5, harga_beli = $6, harga_pasar = $7, tokped_keyword = $8, 
		    keterangan = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10`

	result, err := s.db.Exec(query,
		inv.NamaBarang,
		inv.StokDimiliki,
		inv.StokTerjual,
		inv.StokMasuk,
		inv.HargaJual,
		inv.HargaBeli,
		inv.HargaPasar,
		inv.TokpedKeyword,
		inv.Keterangan,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with ID %d not found", id)
	}

	log.Printf("Updated inventory item ID %d: %s", id, inv.NamaBarang)
	return nil
}

// DeleteInventory deletes an inventory item from PostgreSQL
func (s *InventoryService) DeleteInventory(id int) error {
	query := `DELETE FROM inventory WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with ID %d not found", id)
	}

	log.Printf("Deleted inventory item ID %d", id)
	return nil
}

// InitializeSheet creates headers in the database (not needed for PostgreSQL)
func (s *InventoryService) InitializeSheet() error {
	// Headers are already created in the database schema
	log.Println("Database schema is ready")
	return nil
}

// UpdateMarketPrice updates the market price for a specific item
func (s *InventoryService) UpdateMarketPrice(id int, price int) error {
	query := `UPDATE inventory SET harga_pasar = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := s.db.Exec(query, price, id)
	if err != nil {
		return fmt.Errorf("failed to update market price: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with ID %d not found", id)
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
