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

// GetAllBook retrieves all book items from PostgreSQL
func GetAllBook(bookService *service.BookService) ([]model.Book, error) {
	return bookService.GetAllBook()
}

// DashboardController displays the main book dashboard
func DashboardController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inventories, err := GetAllBook(bookService)
		if err != nil {
			log.Printf("Error fetching book: %v", err)
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
			"len": func(s []model.Book) int {
				return len(s)
			},
			"getTotalStock": func(inventories []model.Book) int {
				total := 0
				for _, inv := range inventories {
					total += inv.Stok
				}
				return total
			},
			"getTotalValue": func(inventories []model.Book) int {
				total := 0
				for _, inv := range inventories {
					total += inv.Stok * inv.Harga
				}
				return total
			},
			"getScrapingCount": func(inventories []model.Book) int {
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

// AddProductController handles adding new book items
func AddProductController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
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
			stok, _ := strconv.Atoi(r.FormValue("stok"))
			harga, _ := strconv.Atoi(r.FormValue("harga"))
			tokpedKeyword := r.FormValue("tokped_keyword")
			keterangan := r.FormValue("keterangan")

			// Simulate market price scraping
			hargaPasar := service.SimulatePriceScraping(tokpedKeyword, harga)

			// Create book model
			book := model.Book{
				NamaBarang:    namaBarang,
				Stok:          stok,
				Harga:         harga,
				HargaPasar:    hargaPasar,
				TokpedKeyword: tokpedKeyword,
				Keterangan:    keterangan,
			}

			err := bookService.AddBook(book)
			if err != nil {
				log.Printf("Error inserting product: %v", err)
				http.Error(w, "Failed to save product", http.StatusInternalServerError)
				return
			}

			log.Printf("New book item added: %s", namaBarang)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// EditBookController displays the edit form for a specific book
func EditBookController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Ambil ID dari query parameter, misalnya /edit?id=5
			idStr := r.URL.Query().Get("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid book ID", http.StatusBadRequest)
				return
			}

			// Ambil data book dari service
			book, err := bookService.GetBookByID(id)
			if err != nil {
				http.Error(w, "Book not found", http.StatusNotFound)
				return
			}

			// Load template
			fp := filepath.Join("view", "editbook.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				log.Printf("Error parsing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Kirim data ke template
			err = tmpl.Execute(w, book)
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

// UpdateBookController handles updating book items
func UpdateBookController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			namaBarang := r.FormValue("nama_barang")
			stok, _ := strconv.Atoi(r.FormValue("stok"))
			terjual, _ := strconv.Atoi(r.FormValue("terjual"))
			harga, _ := strconv.Atoi(r.FormValue("harga"))
			tokpedKeyword := r.FormValue("tokped_keyword")
			keterangan := r.FormValue("keterangan")

			// Re-scrape market price
			hargaPasar := service.SimulatePriceScraping(tokpedKeyword, harga)

			// Create book model
			book := model.Book{
				ID:            id,
				NamaBarang:    namaBarang,
				Stok:          stok,
				Terjual:       terjual,
				Harga:         harga,
				HargaPasar:    hargaPasar,
				TokpedKeyword: tokpedKeyword,
				Keterangan:    keterangan,
			}

			err := bookService.UpdateBook(id, book)
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

// DeleteBookController handles deleting book items
func DeleteBookController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))

			err := bookService.DeleteBook(id)
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
func ManageProductController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Get all book items
			inventories, err := GetAllBook(bookService)
			if err != nil {
				log.Printf("Error fetching book: %v", err)
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
				"len": func(s []model.Book) int {
					return len(s)
				},
				"getTotalStock": func(inventories []model.Book) int {
					total := 0
					for _, inv := range inventories {
						total += inv.Stok
					}
					return total
				},
				"getTotalValue": func(inventories []model.Book) int {
					total := 0
					for _, inv := range inventories {
						total += inv.Stok * inv.Harga
					}
					return total
				},
				"getScrapingCount": func(inventories []model.Book) int {
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
func ScrapePriceController(bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Get all items with tokped_keyword
			inventories, err := bookService.GetAllBook()
			if err != nil {
				log.Printf("Error fetching items for scraping: %v", err)
				http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
				return
			}

			var updatedCount int
			for _, inv := range inventories {
				if inv.TokpedKeyword != "" {
					// Simulate price scraping
					hargaPasar := service.SimulatePriceScraping(inv.TokpedKeyword, inv.Harga)

					// Update the market price
					err = bookService.UpdateMarketPrice(inv.ID, hargaPasar)
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
