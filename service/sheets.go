package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type SheetsService struct {
	service *sheets.Service
	sheetID string
}

// NewSheetsService creates a new Google Sheets service
func NewSheetsService(sheetID string) (*SheetsService, error) {
	ctx := context.Background()

	// Read credentials from environment variable or file
	credentialsJSON := os.Getenv("GOOGLE_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		// Try to read from file
		credFile := "credentials.json"
		if _, err := os.Stat(credFile); err == nil {
			credBytes, err := os.ReadFile(credFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read credentials file: %v", err)
			}
			credentialsJSON = string(credBytes)
		} else {
			// Use service account key for testing (you need to create this)
			return nil, fmt.Errorf("please set GOOGLE_CREDENTIALS_JSON environment variable or place credentials.json file")
		}
	}

	// Parse credentials
	config, err := google.JWTConfigFromJSON([]byte(credentialsJSON), sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	client := config.Client(ctx)

	// Create sheets service
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %v", err)
	}

	return &SheetsService{
		service: service,
		sheetID: sheetID,
	}, nil
}

// GetAllInventory retrieves all inventory items from Google Sheets
func (s *SheetsService) GetAllInventory() ([]model.Inventory, error) {
	rangeStr := "Sheet1!A1:K1000" // Adjust range as needed
	resp, err := s.service.Spreadsheets.Values.Get(s.sheetID, rangeStr).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get data from sheets: %v", err)
	}

	var inventories []model.Inventory
	
	// Skip header row (row 1), start from row 2
	if len(resp.Values) <= 1 {
		return inventories, nil
	}

	for i, row := range resp.Values[1:] { // Skip header
		if len(row) < 9 {
			continue // Skip incomplete rows
		}

		inventory := model.Inventory{
			ID: i + 1, // Row number as ID
		}

		// Parse each field
		if len(row) > 0 && row[0] != nil {
			inventory.NamaBarang = fmt.Sprintf("%v", row[0])
		}
		if len(row) > 1 && row[1] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[1])); err == nil {
				inventory.StokDimiliki = val
			}
		}
		if len(row) > 2 && row[2] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[2])); err == nil {
				inventory.StokTerjual = val
			}
		}
		if len(row) > 3 && row[3] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[3])); err == nil {
				inventory.StokMasuk = val
			}
		}
		if len(row) > 4 && row[4] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[4])); err == nil {
				inventory.HargaJual = val
			}
		}
		if len(row) > 5 && row[5] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[5])); err == nil {
				inventory.HargaBeli = val
			}
		}
		if len(row) > 6 && row[6] != nil {
			if val, err := strconv.Atoi(fmt.Sprintf("%v", row[6])); err == nil {
				inventory.HargaPasar = val
			}
		}
		if len(row) > 7 && row[7] != nil {
			inventory.TokpedKeyword = fmt.Sprintf("%v", row[7])
		}
		if len(row) > 8 && row[8] != nil {
			inventory.Keterangan = fmt.Sprintf("%v", row[8])
		}

		inventories = append(inventories, inventory)
	}

	return inventories, nil
}

// AddInventory adds a new inventory item to Google Sheets
func (s *SheetsService) AddInventory(inv model.Inventory) error {
	// Find the next empty row
	inventories, err := s.GetAllInventory()
	if err != nil {
		return err
	}

	nextRow := len(inventories) + 2 // +2 because we skip header and arrays are 0-indexed

	// Prepare data row
	values := [][]interface{}{
		{
			inv.NamaBarang,
			inv.StokDimiliki,
			inv.StokTerjual,
			inv.StokMasuk,
			inv.HargaJual,
			inv.HargaBeli,
			inv.HargaPasar,
			inv.TokpedKeyword,
			inv.Keterangan,
		},
	}

	rangeStr := fmt.Sprintf("Sheet1!A%d:I%d", nextRow, nextRow)
	
	// Create value range
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err = s.service.Spreadsheets.Values.Update(s.sheetID, rangeStr, valueRange).
		ValueInputOption("RAW").Do()
	
	if err != nil {
		return fmt.Errorf("failed to add inventory: %v", err)
	}

	log.Printf("Added inventory item: %s", inv.NamaBarang)
	return nil
}

// UpdateInventory updates an existing inventory item in Google Sheets
func (s *SheetsService) UpdateInventory(id int, inv model.Inventory) error {
	row := id + 1 // +1 because arrays are 0-indexed but sheets start from 1, and we have header

	values := [][]interface{}{
		{
			inv.NamaBarang,
			inv.StokDimiliki,
			inv.StokTerjual,
			inv.StokMasuk,
			inv.HargaJual,
			inv.HargaBeli,
			inv.HargaPasar,
			inv.TokpedKeyword,
			inv.Keterangan,
		},
	}

	rangeStr := fmt.Sprintf("Sheet1!A%d:I%d", row, row)
	
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err := s.service.Spreadsheets.Values.Update(s.sheetID, rangeStr, valueRange).
		ValueInputOption("RAW").Do()
	
	if err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	log.Printf("Updated inventory item ID %d: %s", id, inv.NamaBarang)
	return nil
}

// DeleteInventory deletes an inventory item from Google Sheets
func (s *SheetsService) DeleteInventory(id int) error {
	row := id + 1 // +1 because arrays are 0-indexed but sheets start from 1, and we have header

	// Clear the row by setting empty values
	values := [][]interface{}{
		{"", "", "", "", "", "", "", "", ""},
	}

	rangeStr := fmt.Sprintf("Sheet1!A%d:I%d", row, row)
	
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err := s.service.Spreadsheets.Values.Update(s.sheetID, rangeStr, valueRange).
		ValueInputOption("RAW").Do()
	
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %v", err)
	}

	log.Printf("Deleted inventory item ID %d", id)
	return nil
}

// InitializeSheet creates headers in the Google Sheet if they don't exist
func (s *SheetsService) InitializeSheet() error {
	// Check if headers exist
	rangeStr := "Sheet1!A1:I1"
	resp, err := s.service.Spreadsheets.Values.Get(s.sheetID, rangeStr).Do()
	if err != nil {
		return fmt.Errorf("failed to check headers: %v", err)
	}

	// If headers don't exist or are empty, create them
	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		headers := [][]interface{}{
			{
				"Nama Barang",
				"Stok Dimiliki", 
				"Stok Terjual",
				"Stok Masuk",
				"Harga Jual",
				"Harga Beli",
				"Harga Pasar",
				"Tokopedia Keyword",
				"Keterangan",
			},
		}

		valueRange := &sheets.ValueRange{
			Values: headers,
		}

		_, err = s.service.Spreadsheets.Values.Update(s.sheetID, rangeStr, valueRange).
			ValueInputOption("RAW").Do()
		
		if err != nil {
			return fmt.Errorf("failed to create headers: %v", err)
		}

		log.Println("Created headers in Google Sheet")
	}

	return nil
}

// UpdateMarketPrice updates the market price for a specific row
func (s *SheetsService) UpdateMarketPrice(row int, price int) error {
	rangeStr := fmt.Sprintf("Sheet1!G%d", row+2) // +2 for header and 0-indexing
	
	values := [][]interface{}{
		{price},
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err := s.service.Spreadsheets.Values.Update(s.sheetID, rangeStr, valueRange).
		ValueInputOption("RAW").Do()
	
	if err != nil {
		return fmt.Errorf("failed to update market price: %v", err)
	}

	return nil
}

// simulatePriceScraping simulates price scraping from Tokopedia
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

// Helper function to create a mock sheets service for testing without Google API
func CreateMockSheetsService() *MockSheetsService {
	return &MockSheetsService{
		data: make([]model.Inventory, 0),
	}
}

// MockSheetsService for testing without Google API
type MockSheetsService struct {
	data []model.Inventory
}

func (m *MockSheetsService) GetAllInventory() ([]model.Inventory, error) {
	return m.data, nil
}

func (m *MockSheetsService) AddInventory(inv model.Inventory) error {
	inv.ID = len(m.data) + 1
	m.data = append(m.data, inv)
	return nil
}

func (m *MockSheetsService) UpdateInventory(id int, inv model.Inventory) error {
	for i := range m.data {
		if m.data[i].ID == id {
			inv.ID = id
			m.data[i] = inv
			return nil
		}
	}
	return fmt.Errorf("inventory with ID %d not found", id)
}

func (m *MockSheetsService) DeleteInventory(id int) error {
	for i := range m.data {
		if m.data[i].ID == id {
			m.data = append(m.data[:i], m.data[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("inventory with ID %d not found", id)
}

func (m *MockSheetsService) InitializeSheet() error {
	return nil
}

func (m *MockSheetsService) UpdateMarketPrice(row int, price int) error {
	if row < len(m.data) {
		m.data[row].HargaPasar = price
	}
	return nil
}
