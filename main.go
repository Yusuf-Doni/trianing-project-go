package main

import (
	"log"
	"net/http"

	"github.com/Yusuf-Doni/web-go-CRUD/database"
	"github.com/Yusuf-Doni/web-go-CRUD/routes"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

func main() {
	// Initialize PostgreSQL database connection
	log.Println("Initializing PostgreSQL database connection...")
	db := database.InitDatabase()
	defer db.Close()

	// Create inventory service
	inventoryService := service.NewInventoryService(db)

	// Initialize database schema
	err := inventoryService.InitializeSheet()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	server := http.NewServeMux()

	routes.MapRoutes(server, inventoryService, db)

	log.Println("🚀 Server starting on :8000")
	log.Println("🔐 Login: http://localhost:8000/login")
	log.Println("📝 Register: http://localhost:8000/register")
	log.Println("📊 Dashboard: http://localhost:8000/dashboard")
	log.Println("➕ Add Product: http://localhost:8000/addproduct")
	log.Println("👤 Default Admin: username=admin, password=admin123")
	
	http.ListenAndServe(":8000", server)
}
