package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

// GetAllInventory retrieves all inventory items from PostgreSQL
func GetAllInventory(inventoryService *service.InventoryService) ([]model.Inventory, error) {
	return inventoryService.GetAllInventory()
}

// DashboardController displays the main inventory dashboard
func DashboardController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inventories, err := GetAllInventory(inventoryService)
		if err != nil {
			log.Printf("Error fetching inventory: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fp := filepath.Join("view", "index.html")

		// Create template with custom functions
		tmpl := template.New("index.html").Funcs(template.FuncMap{
			"formatNumber": func(n int) string {
				return formatCurrency(n)
			},
			"sub": func(a, b int) int {
				return a - b
			},
			"len": func(s []model.Inventory) int {
				return len(s)
			},
			"getTotalStock": func(inventories []model.Inventory) int {
				total := 0
				for _, inv := range inventories {
					total += inv.StokDimiliki
				}
				return total
			},
			"getTotalValue": func(inventories []model.Inventory) int {
				total := 0
				for _, inv := range inventories {
					total += inv.StokDimiliki * inv.HargaJual
				}
				return total
			},
			"getScrapingCount": func(inventories []model.Inventory) int {
				count := 0
				for _, inv := range inventories {
					if inv.TokpedKeyword != "" {
						count++
					}
				}
				return count
			},
		})

		tmpl, err = tmpl.ParseFiles(fp)
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, inventories)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// formatCurrency formats number as Indonesian currency
func formatCurrency(amount int) string {
	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return str
	}

	result := ""
	for i, digit := range str {
		if (n-i)%3 == 0 && i > 0 {
			result += "."
		}
		result += string(digit)
	}
	return result
}

// AddProductController handles adding new inventory items
func AddProductController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fp := filepath.Join("view", "addproduct.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if r.Method == "POST" {
			r.ParseForm()

			namaBarang := r.FormValue("nama_barang")
			stokDimiliki, _ := strconv.Atoi(r.FormValue("stok_dimiliki"))
			stokTerjual, _ := strconv.Atoi(r.FormValue("stok_terjual"))
			stokMasuk, _ := strconv.Atoi(r.FormValue("stok_masuk"))
			hargaJual, _ := strconv.Atoi(r.FormValue("harga_jual"))
			hargaBeli, _ := strconv.Atoi(r.FormValue("harga_beli"))
			tokpedKeyword := r.FormValue("tokped_keyword")
			keterangan := r.FormValue("keterangan")

			// Simulate market price scraping
			hargaPasar := service.SimulatePriceScraping(tokpedKeyword, hargaJual)

			// Create inventory model
			inventory := model.Inventory{
				NamaBarang:    namaBarang,
				StokDimiliki:  stokDimiliki,
				StokTerjual:   stokTerjual,
				StokMasuk:     stokMasuk,
				HargaJual:     hargaJual,
				HargaBeli:     hargaBeli,
				HargaPasar:    hargaPasar,
				TokpedKeyword: tokpedKeyword,
				Keterangan:    keterangan,
			}

			err := inventoryService.AddInventory(inventory)
			if err != nil {
				log.Printf("Error inserting product: %v", err)
				http.Error(w, "Failed to save product", http.StatusInternalServerError)
				return
			}

			log.Printf("New inventory item added: %s", namaBarang)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// UpdateInventoryController handles updating inventory items
func UpdateInventoryController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			namaBarang := r.FormValue("nama_barang")
			stokDimiliki, _ := strconv.Atoi(r.FormValue("stok_dimiliki"))
			stokTerjual, _ := strconv.Atoi(r.FormValue("stok_terjual"))
			stokMasuk, _ := strconv.Atoi(r.FormValue("stok_masuk"))
			hargaJual, _ := strconv.Atoi(r.FormValue("harga_jual"))
			hargaBeli, _ := strconv.Atoi(r.FormValue("harga_beli"))
			tokpedKeyword := r.FormValue("tokped_keyword")
			keterangan := r.FormValue("keterangan")

			// Re-scrape market price
			hargaPasar := service.SimulatePriceScraping(tokpedKeyword, hargaJual)

			// Create inventory model
			inventory := model.Inventory{
				ID:            id,
				NamaBarang:    namaBarang,
				StokDimiliki:  stokDimiliki,
				StokTerjual:   stokTerjual,
				StokMasuk:     stokMasuk,
				HargaJual:     hargaJual,
				HargaBeli:     hargaBeli,
				HargaPasar:    hargaPasar,
				TokpedKeyword: tokpedKeyword,
				Keterangan:    keterangan,
			}

			err := inventoryService.UpdateInventory(id, inventory)
			if err != nil {
				log.Printf("Error updating product: %v", err)
				http.Error(w, "Failed to update product", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// DeleteInventoryController handles deleting inventory items
func DeleteInventoryController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))

			err := inventoryService.DeleteInventory(id)
			if err != nil {
				log.Printf("Error deleting product: %v", err)
				http.Error(w, "Failed to delete product", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ManageProductController displays the product management page
func ManageProductController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Get all inventory items
			inventories, err := GetAllInventory(inventoryService)
			if err != nil {
				log.Printf("Error fetching inventory: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			fp := filepath.Join("view", "manageproduk.html")

			// Create template with custom functions
			tmpl := template.New("manageproduk.html").Funcs(template.FuncMap{
				"formatNumber": func(n int) string {
					return formatCurrency(n)
				},
				"sub": func(a, b int) int {
					return a - b
				},
				"len": func(s []model.Inventory) int {
					return len(s)
				},
				"getTotalStock": func(inventories []model.Inventory) int {
					total := 0
					for _, inv := range inventories {
						total += inv.StokDimiliki
					}
					return total
				},
				"getTotalValue": func(inventories []model.Inventory) int {
					total := 0
					for _, inv := range inventories {
						total += inv.StokDimiliki * inv.HargaJual
					}
					return total
				},
				"getScrapingCount": func(inventories []model.Inventory) int {
					count := 0
					for _, inv := range inventories {
						if inv.TokpedKeyword != "" {
							count++
						}
					}
					return count
				},
			})

			tmpl, err = tmpl.ParseFiles(fp)
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, inventories)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ScrapePriceController triggers price scraping for all items
func ScrapePriceController(inventoryService *service.InventoryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Get all items with tokped_keyword
			inventories, err := inventoryService.GetAllInventory()
			if err != nil {
				log.Printf("Error fetching items for scraping: %v", err)
				http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
				return
			}

			var updatedCount int
			for _, inv := range inventories {
				if inv.TokpedKeyword != "" {
					// Simulate price scraping
					hargaPasar := service.SimulatePriceScraping(inv.TokpedKeyword, inv.HargaJual)

					// Update the market price
					err = inventoryService.UpdateMarketPrice(inv.ID, hargaPasar)
					if err == nil {
						updatedCount++
					}
				}
			}

			response := map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("Successfully updated %d items", updatedCount),
				"count":   updatedCount,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
