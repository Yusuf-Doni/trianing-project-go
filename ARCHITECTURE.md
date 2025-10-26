# Arsitektur Sistem Book Management

## Overview
Sistem book management dengan PostgreSQL database menggunakan arsitektur monolitik. Sistem ini fokus pada manajemen book dengan simulasi scraping harga pasar.

## Komponen Sistem

```
┌─────────────────────────────────────────────────────────────┐
│                    FRONTEND (HTML + Bootstrap)              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Dashboard     │  │  Add Product    │  │   Actions    │ │
│  │   - Tabel Data  │  │  - Form Input   │  │   - Scraping │ │
│  │   - Real-time   │  │  - Validation   │  │   - CRUD     │ │
│  │   - Analytics   │  │  - AI Features  │  │   - Refresh  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                 BACKEND (Golang)                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Controllers   │  │     Models      │  │    Routes    │ │
│  │   - CRUD Ops    │  │   - Book   │  │   - REST API │ │
│  │   - Validation  │  │   - Data Types  │  │   - Handlers │ │
│  │   - Business    │  │   - Structures  │  │   - Mapping  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                DATABASE (PostgreSQL)                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Book     │  │   Auto Schema   │  │  Sample Data │ │
│  │   - CRUD Data   │  │   - Migration   │  │   - Testing  │ │
│  │   - Relations   │  │   - Indexes     │  │   - Seeds    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│            EXTERNAL SERVICES (Simulated)                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Tokopedia     │  │   Gemini AI     │  │   Scraping   │ │
│  │   - Price Data  │  │   - Validation  │  │   - Keywords │ │
│  │   - Search API  │  │   - Normalize   │  │   - Results  │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Input Process
```
User Input → Form Validation → Business Logic → Database Storage
```

### 2. Scraping Process
```
Keyword → Tokopedia Search → AI Validation → Price Update → Database
```

### 3. Display Process
```
Database Query → Data Processing → Template Rendering → HTML Output
```

## Database Schema

### Tabel: book
```sql
CREATE TABLE book (
    id SERIAL PRIMARY KEY,
    nama_barang TEXT NOT NULL,           -- Kolom A: Nama produk yang dinormalisasi
    stok_dimiliki INTEGER DEFAULT 0,     -- Kolom B: Stok fisik di gudang
    stok_terjual INTEGER DEFAULT 0,      -- Kolom C: Jumlah terjual
    stok_masuk INTEGER DEFAULT 0,        -- Kolom D: Jumlah masuk terakhir
    harga_jual INTEGER DEFAULT 0,        -- Kolom E: Harga jual yang ditetapkan
    harga_beli INTEGER DEFAULT 0,        -- Kolom F: Harga modal
    harga_pasar INTEGER DEFAULT 0,       -- Kolom G: Harga pasar (hasil scraping)
    tokped_keyword TEXT,                 -- Kolom H: Keyword untuk scraping
    keterangan TEXT,                     -- Kolom I: Deskripsi tambahan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

### REST API Structure
```
GET    /dashboard           → Dashboard utama dengan data book
GET    /addproduct          → Form tambah produk baru
POST   /addproduct          → Submit form tambah produk
POST   /update              → Update data produk existing
POST   /delete              → Hapus data produk
POST   /scrape              → Trigger scraping harga Tokopedia
GET    /                    → Hello World endpoint
```

## Fitur Simulasi Scraping

### Simulasi Scraping Harga Pasar
1. **Input**: Keyword dari field `tokped_keyword`
2. **Process**: Simulasi algoritma pencarian harga
3. **Validation**: Validasi harga masuk akal (80%-120% dari harga jual)
4. **Output**: Harga pasar yang dinormalisasi

### Validasi Harga
1. **Price Range Check**: Memastikan harga masuk akal (50%-150% dari harga jual)
2. **Product Name Normalization**: Standardisasi nama produk
3. **Data Quality**: Filter hasil yang tidak valid
4. **Simulation Algorithm**: Menggunakan variasi berdasarkan panjang keyword

## Security & Performance

### Security Measures
- Input validation pada client dan server
- SQL injection protection dengan prepared statements
- CSRF protection pada form submissions
- Error handling yang proper tanpa informasi sensitif

### Performance Optimizations
- Database indexing pada field yang sering di-query
- Template caching untuk rendering HTML
- Connection pooling untuk database
- Lazy loading untuk data besar

## Deployment Architecture

### Development Environment
```
Local Machine → PostgreSQL → Go Application → Browser
```

### Production Environment
```
Load Balancer → Go Application → PostgreSQL Cluster
```

## Monitoring & Logging

### Application Metrics
- Response time untuk setiap endpoint
- Database query performance
- Scraping success rate
- Error rate dan exception handling

### Business Metrics
- Total produk dalam book
- Nilai total book
- Harga pasar vs harga jual comparison
- Scraping accuracy rate

## Future Enhancements

### Planned Features
1. **Real Tokopedia API Integration**
2. **Advanced Analytics Dashboard**
3. **Automated Price Alerts**
4. **Multi-vendor Support**
5. **Mobile App Integration**
6. **Real-time Price Monitoring**

### Scalability Considerations
1. **Microservices Architecture**: Memisahkan scraping service
2. **Message Queue**: Untuk async processing
3. **Caching Layer**: Redis untuk data yang sering diakses
4. **Database Sharding**: Untuk data yang sangat besar
5. **API Gateway**: Untuk rate limiting dan monitoring

---

**Catatan**: Sistem ini menggunakan PostgreSQL sebagai database utama dan simulasi untuk fitur scraping harga pasar sesuai dengan kebutuhan development dan testing.
