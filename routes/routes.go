package routes

import (
	"database/sql"
	"net/http"

	"github.com/Yusuf-Doni/web-go-CRUD/controller"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

func MapRoutes(server *http.ServeMux, bookService *service.BookService, db *sql.DB) {
	// Inisialisasi service baru untuk cart
	cartService := service.NewCartService(db)

	// Public routes (no authentication required)
	server.HandleFunc("/login", controller.LoginController(db))
	server.HandleFunc("/register", controller.RegisterController(db))
	server.HandleFunc("/logout", controller.LogoutController(db))

	// Protected routes (authentication required)
	server.HandleFunc("/dashboard", controller.RequireAuth(db, controller.DashboardController(bookService)))
	server.HandleFunc("/addproduct", controller.RequireAuth(db, controller.AddProductController(bookService)))
	server.HandleFunc("/manageproduk", controller.RequireAuth(db, controller.ManageProductController(bookService)))
	server.HandleFunc("/edit", controller.RequireAuth(db, controller.EditBookController(bookService)))
	server.HandleFunc("/update", controller.RequireAuth(db, controller.UpdateBookController(bookService)))
	server.HandleFunc("/delete", controller.RequireAuth(db, controller.DeleteBookController(bookService)))
	server.HandleFunc("/scrape", controller.RequireAuth(db, controller.ScrapePriceController(bookService)))

	// ================================
	// 🛒 CART FEATURE (new routes buat Cart)
	// ================================
	server.HandleFunc("/cart/add", controller.RequireAuth(db, controller.AddToCartController(cartService, bookService)))
	server.HandleFunc("/cart/get", controller.RequireAuth(db, controller.GetCartController(cartService)))
	server.HandleFunc("/cart/remove", controller.RequireAuth(db, controller.RemoveFromCartController(cartService)))

	// Root route - redirect to login or dashboard based on auth status
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if controller.IsLoggedIn(r) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
}
