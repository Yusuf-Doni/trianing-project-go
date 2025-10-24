package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
)

type AppsScriptService struct {
	webAppURL string
	client    *http.Client
}

// NewAppsScriptService creates a new Google Apps Script service
func NewAppsScriptService(webAppURL string) *AppsScriptService {
	return &AppsScriptService{
		webAppURL: webAppURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAllInventory retrieves all inventory items from Google Apps Script
func (s *AppsScriptService) GetAllInventory() ([]model.Inventory, error) {
	resp, err := s.client.Get(s.webAppURL + "?action=getInventory")
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Status string             `json:"status"`
		Data   []model.Inventory `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("API returned error status")
	}

	return result.Data, nil
}

// AddInventory adds a new inventory item via Google Apps Script
func (s *AppsScriptService) AddInventory(inv model.Inventory) error {
	data := map[string]interface{}{
		"action":        "addInventory",
		"nama_barang":   inv.NamaBarang,
		"stok_dimiliki": inv.StokDimiliki,
		"stok_terjual":  inv.StokTerjual,
		"stok_masuk":    inv.StokMasuk,
		"harga_jual":    inv.HargaJual,
		"harga_beli":    inv.HargaBeli,
		"harga_pasar":   inv.HargaPasar,
		"tokped_keyword": inv.TokpedKeyword,
		"keterangan":    inv.Keterangan,
	}

	return s.sendRequest(data)
}

// UpdateInventory updates an existing inventory item via Google Apps Script
func (s *AppsScriptService) UpdateInventory(id int, inv model.Inventory) error {
	data := map[string]interface{}{
		"action":        "updateInventory",
		"id":            id,
		"nama_barang":   inv.NamaBarang,
		"stok_dimiliki": inv.StokDimiliki,
		"stok_terjual":  inv.StokTerjual,
		"stok_masuk":    inv.StokMasuk,
		"harga_jual":    inv.HargaJual,
		"harga_beli":    inv.HargaBeli,
		"harga_pasar":   inv.HargaPasar,
		"tokped_keyword": inv.TokpedKeyword,
		"keterangan":    inv.Keterangan,
	}

	return s.sendRequest(data)
}

// DeleteInventory deletes an inventory item via Google Apps Script
func (s *AppsScriptService) DeleteInventory(id int) error {
	data := map[string]interface{}{
		"action": "deleteInventory",
		"id":     id,
	}

	return s.sendRequest(data)
}

// InitializeSheet creates headers in the Google Sheet (not needed for Apps Script)
func (s *AppsScriptService) InitializeSheet() error {
	// Headers are already created in the Google Sheet
	return nil
}

// UpdateMarketPrice updates the market price for a specific item
func (s *AppsScriptService) UpdateMarketPrice(row int, price int) error {
	// This would need to be implemented in the Apps Script
	// For now, we'll use the general update method
	return nil
}

// ScrapePrice triggers price scraping via Google Apps Script
func (s *AppsScriptService) ScrapePrice(keyword string, basePrice int) (int, error) {
	data := map[string]interface{}{
		"action":        "scrapePrice",
		"tokped_keyword": keyword,
		"harga_jual":    basePrice,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal data: %v", err)
	}

	resp, err := s.client.Post(s.webAppURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Status     string `json:"status"`
		HargaPasar int    `json:"harga_pasar"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Status != "success" {
		return 0, fmt.Errorf("API returned error status")
	}

	return result.HargaPasar, nil
}

// sendRequest sends a POST request to Google Apps Script
func (s *AppsScriptService) sendRequest(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	resp, err := s.client.Post(s.webAppURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Status != "success" {
		return fmt.Errorf("API error: %s", result.Message)
	}

	return nil
}

// SimulatePriceScraping simulates price scraping (for compatibility)
func (s *AppsScriptService) SimulatePriceScraping(keyword string, basePrice int) int {
	// Use the Apps Script scraping function
	price, err := s.ScrapePrice(keyword, basePrice)
	if err != nil {
		// Fallback to local simulation
		return SimulatePriceScraping(keyword, basePrice)
	}
	return price
}
