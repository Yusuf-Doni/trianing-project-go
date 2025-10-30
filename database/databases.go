package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase inisialisasi koneksi ke PostgreSQL
func InitDatabase() *sql.DB {
	// Ambil variabel environment
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Susun DSN PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, sslmode)

	// Buat koneksi
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("gagal open DB: %v", err)
	}

	// Cek koneksi
	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal connect ke DB: %v", err)
	}

	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		role VARCHAR(20) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createUsersTableSQL); err != nil {
		log.Fatalf("Gagal membuat tabel users: %v", err)
	}

	createSessionsTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(64) PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		username VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,	
		expires_at TIMESTAMP NOT NULL
	);`

	if _, err = DB.Exec(createSessionsTableSQL); err != nil {
		log.Fatalf("Gagal membuat tabel sessions: %v", err)
	}

	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	`
	if _, err = DB.Exec(createIndexes); err != nil {
		log.Printf("Warning: gagal membuat index: %v", err)
	}

	createBooksTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		nama_barang TEXT NOT NULL,
		stok INTEGER DEFAULT 0,
		terjual INTEGER DEFAULT 0,
		harga INTEGER DEFAULT 0,
		harga_pasar INTEGER DEFAULT 0,
		tokped_keyword TEXT,
		keterangan TEXT,
		gambar_buku TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createBooksTableSQL); err != nil {
		log.Fatalf("Gagal membuat tabel books: %v", err)
	}

	createCartsTableSQL := `
	CREATE TABLE IF NOT EXISTS carts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
		jumlah INTEGER NOT NULL DEFAULT 1,
		harga INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createCartsTableSQL); err != nil {
		log.Fatalf("Gagal membuat tabel carts: %v", err)
	}

	createOrdersTableSQL := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		total_harga INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createOrdersTableSQL); err != nil {
		log.Fatalf("gagal membuat tabel orders: %v", err)
	}

	createOrderItemsTableSQL := `
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
		jumlah INTEGER NOT NULL DEFAULT 1,
		harga INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err = DB.Exec(createOrderItemsTableSQL); err != nil {
		log.Fatalf("Gagal membuat tabel order_items: %v", err)
	}

	createPaymentTable := `
CREATE TABLE IF NOT EXISTS payment (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    total_item INTEGER NOT NULL DEFAULT 0,
    total_harga INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0, -- 0=draft, 1=sukses, 2=batal
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

	_, err = DB.Exec(createPaymentTable)
	if err != nil {
		log.Fatalf("gagal membuat tabel payment: %v", err)
	}

	createPaymentDetailTable := `
CREATE TABLE IF NOT EXISTS payment_detail (
    id SERIAL PRIMARY KEY,
    payment_id INTEGER NOT NULL,
    book_id INTEGER NOT NULL,
    jumlah INTEGER NOT NULL DEFAULT 1,
    harga_satuan INTEGER NOT NULL DEFAULT 0,
    total_harga INTEGER NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0, -- 0=draft, 1=sukses, 2=batal
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (payment_id) REFERENCES payment(id) ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);`

	_, err = DB.Exec(createPaymentDetailTable)
	if err != nil {
		log.Fatalf("gagal membuat tabel payment_detail: %v", err)
	}

	// Insert default admin user jika belum ada
	insertAdminSQL := `
	INSERT INTO users (username, password, email, role) VALUES 
	('admin', 'admin123', 'admin@book.com', 'admin')
	ON CONFLICT (username) DO NOTHING;`

	_, err = DB.Exec(insertAdminSQL)
	if err != nil {
		log.Printf("Warning: gagal insert admin user: %v", err)
	}

	log.Println("Koneksi ke PostgreSQL berhasil dan tabel book siap!")
	return DB
}
