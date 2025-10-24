@echo off
echo ========================================
echo   Sistem Manajemen Inventory
echo   Google Sheets Integration
echo ========================================
echo.

echo [INFO] Menjalankan aplikasi dengan Google Sheets...
echo [INFO] Dashboard: http://localhost:8000/dashboard
echo [INFO] Tambah Produk: http://localhost:8000/addproduct
echo.
echo [NOTE] Jika Google Sheets tidak bisa diakses, aplikasi akan menggunakan mock service
echo [NOTE] Lihat setup_google_sheets.md untuk konfigurasi Google Sheets
echo.

go run main.go
