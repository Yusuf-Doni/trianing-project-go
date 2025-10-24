package routes

import (
	"database/sql"
	"net/http"

	"github.com/Yusuf-Doni/web-go-CRUD/controller"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

func MapRoutes(server *http.ServeMux, inventoryService *service.InventoryService, db *sql.DB) {
	// Public routes (no authentication required)
	server.HandleFunc("/login", controller.LoginController(db))
	server.HandleFunc("/register", controller.RegisterController(db))
	server.HandleFunc("/logout", controller.LogoutController(db))
	
	// Protected routes (authentication required)
	server.HandleFunc("/dashboard", controller.RequireAuth(db, controller.DashboardController(inventoryService)))
	server.HandleFunc("/addproduct", controller.RequireAuth(db, controller.AddProductController(inventoryService)))
	server.HandleFunc("/manageproduk", controller.RequireAuth(db, controller.ManageProductController(inventoryService)))
	server.HandleFunc("/update", controller.RequireAuth(db, controller.UpdateInventoryController(inventoryService)))
	server.HandleFunc("/delete", controller.RequireAuth(db, controller.DeleteInventoryController(inventoryService)))
	server.HandleFunc("/scrape", controller.RequireAuth(db, controller.ScrapePriceController(inventoryService)))
	
	// Root route - redirect to login or dashboard based on auth status
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if user is logged in
		if controller.IsLoggedIn(r) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
}
