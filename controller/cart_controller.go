package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Yusuf-Doni/web-go-CRUD/model"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

// AddToCartController menangani request untuk menambahkan item ke cart
func AddToCartController(cartService *service.CartService, bookService *service.BookService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Printf("Gagal parsing form: %v", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		bookID, _ := strconv.Atoi(r.FormValue("book_id"))
		jumlah, _ := strconv.Atoi(r.FormValue("jumlah"))
		if jumlah <= 0 {
			jumlah = 1
		}

		book, err := bookService.GetBookByID(bookID)
		if err != nil {
			log.Printf("Gagal mendapatkan data buku: %v", err)
			http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
			return
		}

		cart := model.Cart{
			UserID: userID,
			BookID: bookID,
			Jumlah: jumlah,
			Harga:  book.Harga * jumlah,
		}

		err = cartService.AddToCart(cart)
		if err != nil {
			log.Printf("Gagal menambahkan ke cart: %v", err)
			http.Error(w, "Gagal menambahkan ke cart", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Produk berhasil ditambahkan ke cart",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetCartController menampilkan semua item di cart user
func GetCartController(cartService *service.CartService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil || userID == 0 {
			http.Error(w, "User ID tidak valid", http.StatusBadRequest)
			return
		}

		carts, err := cartService.GetCartByUser(userID)
		if err != nil {
			log.Printf("Gagal mengambil cart: %v", err)
			http.Error(w, "Gagal mengambil cart", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(carts)
	}
}

// RemoveFromCartController mengurangi jumlah item atau menghapus dari cart
func RemoveFromCartController(cartService *service.CartService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		cartID, err := strconv.Atoi(r.FormValue("cart_id"))
		if err != nil || cartID == 0 {
			http.Error(w, "Cart ID tidak valid", http.StatusBadRequest)
			return
		}

		newJumlah, err := strconv.Atoi(r.FormValue("jumlah"))
		if err != nil {
			http.Error(w, "Jumlah tidak valid", http.StatusBadRequest)
			return
		}

		err = cartService.UpdateCartQuantity(cartID, newJumlah)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Cart berhasil diperbarui",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
