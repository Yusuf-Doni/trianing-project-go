package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase inisialisasi koneksi ke PostgreSQL
func InitDatabase() *sql.DB {
	// Format DSN: "postgres://username:password@host:port/dbname"
	dsn := "postgres://postgres:password@localhost:5434/postgres?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("gagal open DB: %v", err)
	}

	// cek koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatalf("gagal connect DB: %v", err)
	}

	// Buat tabel inventory sesuai struktur Google Sheets
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		nama_barang TEXT NOT NULL,
		stok_dimiliki INTEGER DEFAULT 0,
		stok_terjual INTEGER DEFAULT 0,
		stok_masuk INTEGER DEFAULT 0,
		harga_jual INTEGER DEFAULT 0,
		harga_beli INTEGER DEFAULT 0,
		harga_pasar INTEGER DEFAULT 0,
		tokped_keyword TEXT,
		keterangan TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("gagal membuat tabel: %v", err)
	}

	log.Println("✅ Koneksi ke PostgreSQL berhasil dan tabel inventory siap!")
	return DB
}
