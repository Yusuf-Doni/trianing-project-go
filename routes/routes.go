package routes

import (
	"net/http"

	"github.com/Yusuf-Doni/web-go-CRUD/controller"
	"github.com/Yusuf-Doni/web-go-CRUD/service"
)

func MapRoutes(server *http.ServeMux, sheetsService service.SheetsServiceInterface) {
	server.HandleFunc("/", controller.HelloWorldController())
	server.HandleFunc("/dashboard", controller.DashboardController(sheetsService))
	server.HandleFunc("/addproduct", controller.AddProductController(sheetsService))
	server.HandleFunc("/update", controller.UpdateInventoryController(sheetsService))
	server.HandleFunc("/delete", controller.DeleteInventoryController(sheetsService))
	server.HandleFunc("/scrape", controller.ScrapePriceController(sheetsService))
}
