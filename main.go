package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Yusuf-Doni/web-go-CRUD/routes"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

func main() {
	// Get Google Apps Script Web App URL from environment variable
	webAppURL := os.Getenv("GOOGLE_APPS_SCRIPT_URL")
	if webAppURL == "" {
		log.Println("No Google Apps Script URL provided. Using mock service for testing...")
		log.Println("Set GOOGLE_APPS_SCRIPT_URL environment variable to use Google Apps Script.")
	}

	var sheetsService service.SheetsServiceInterface

	// Try Google Apps Script first
	if webAppURL != "" {
		log.Println("Initializing Google Apps Script service...")
		sheetsService = service.NewAppsScriptService(webAppURL)
		
		// Test connection
		_, err := sheetsService.GetAllInventory()
		if err != nil {
			log.Printf("Failed to connect to Google Apps Script: %v", err)
			log.Println("Falling back to mock service...")
			sheetsService = service.CreateMockSheetsService()
		} else {
			log.Println("✅ Successfully connected to Google Apps Script!")
		}
	} else {
		// Fallback to mock service
		log.Println("Using mock service for testing...")
		sheetsService = service.CreateMockSheetsService()
	}

	// Initialize sheet headers
	err := sheetsService.InitializeSheet()
	if err != nil {
		log.Printf("Warning: Failed to initialize sheet headers: %v", err)
	}

	server := http.NewServeMux()

	routes.MapRoutes(server, sheetsService)

	log.Println("🚀 Server starting on :8000")
	log.Println("📊 Dashboard: http://localhost:8000/dashboard")
	log.Println("➕ Add Product: http://localhost:8000/addproduct")
	
	http.ListenAndServe(":8000", server)
}
