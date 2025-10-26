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
	dsn := "postgres://postgres:postgres123@localhost:5432/testdb?sslmode=disable"

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

	// Buat tabel book sesuai struktur Google Sheets
	createBookTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		nama_barang TEXT NOT NULL,
		stok INTEGER DEFAULT 0,
		terjual INTEGER DEFAULT 0,
		harga INTEGER DEFAULT 0,
		harga_pasar INTEGER DEFAULT 0,
		tokped_keyword TEXT,
		keterangan TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createBookTableSQL)
	if err != nil {
		log.Fatalf("gagal membuat tabel book: %v", err)
	}

	// Buat tabel users untuk authentication
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

	_, err = DB.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatalf("gagal membuat tabel users: %v", err)
	}

	// Pastikan tabel users sudah dibuat sebelum membuat sessions
	_, err = DB.Exec("SELECT 1 FROM users LIMIT 1")
	if err != nil {
		log.Fatalf("tabel users belum tersedia: %v", err)
	}

	// Buat tabel sessions untuk session management (tanpa foreign key dulu)
	createSessionsTableSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(64) PRIMARY KEY,
		user_id INTEGER NOT NULL,
		username VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,	
		expires_at TIMESTAMP NOT NULL
	);`

	_, err = DB.Exec(createSessionsTableSQL)
	if err != nil {
		log.Fatalf("gagal membuat tabel sessions: %v", err)
	}

	// Tambahkan foreign key constraint setelah tabel dibuat
	addForeignKeySQL := `
	DO $$ 
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE constraint_name = 'sessions_user_id_fkey' 
			AND table_name = 'sessions'
		) THEN
			ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey 
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
		END IF;
	END $$;`

	_, err = DB.Exec(addForeignKeySQL)
	if err != nil {
		log.Printf("Warning: gagal menambahkan foreign key constraint: %v", err)
	}

	// Buat index untuk performa yang lebih baik
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);`

	_, err = DB.Exec(createIndexSQL)
	if err != nil {
		log.Printf("Warning: gagal membuat index: %v", err)
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
