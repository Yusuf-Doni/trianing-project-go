package service

import "github.com/Yusuf-Doni/web-go-CRUD/model"

// SheetsServiceInterface defines the interface for sheets operations
type SheetsServiceInterface interface {
	GetAllInventory() ([]model.Inventory, error)
	AddInventory(inv model.Inventory) error
	UpdateInventory(id int, inv model.Inventory) error
	DeleteInventory(id int) error
	InitializeSheet() error
	UpdateMarketPrice(row int, price int) error
}
