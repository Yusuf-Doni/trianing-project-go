# Setup Google Sheets Integration

## Langkah-langkah Setup Google Sheets API

### 1. Buat Google Cloud Project

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih project yang sudah ada
3. Aktifkan Google Sheets API:
   - Pergi ke "APIs & Services" > "Library"
   - Cari "Google Sheets API"
   - Klik "Enable"

### 2. Buat Service Account

1. Pergi ke "APIs & Services" > "Credentials"
2. Klik "Create Credentials" > "Service Account"
3. Isi nama service account (contoh: `inventory-sheets-service`)
4. Klik "Create and Continue"
5. Pilih role "Editor" atau "Owner"
6. Klik "Done"

### 3. Generate Service Account Key

1. Klik pada service account yang baru dibuat
2. Pergi ke tab "Keys"
3. Klik "Add Key" > "Create new key"
4. Pilih "JSON" format
5. Download file JSON dan simpan sebagai `credentials.json`

### 4. Setup Google Sheet

1. Buat Google Sheet baru
2. Beri nama sheet (contoh: "Inventory Management")
3. Copy Sheet ID dari URL (bagian setelah `/d/` dan sebelum `/edit`)
   ```
   https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
   ```
   Sheet ID: `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms`

### 5. Share Sheet dengan Service Account

1. Di Google Sheet, klik "Share"
2. Tambahkan email service account (dari file credentials.json)
3. Berikan permission "Editor"
4. Klik "Send"

### 6. Konfigurasi Environment Variables

Buat file `.env` atau set environment variables:

```bash
# Windows (PowerShell)
$env:GOOGLE_SHEET_ID="your_sheet_id_here"
$env:GOOGLE_CREDENTIALS_JSON=$(Get-Content credentials.json -Raw)

# Linux/Mac
export GOOGLE_SHEET_ID="your_sheet_id_here"
export GOOGLE_CREDENTIALS_JSON='{"type":"service_account",...}'
```

### 7. Alternatif: File Credentials

Jika tidak ingin menggunakan environment variable, letakkan file `credentials.json` di root project.

## Struktur Google Sheet

Sheet akan memiliki header berikut:

| A | B | C | D | E | F | G | H | I |
|---|---|---|---|---|---|---|---|---|
| Nama Barang | Stok Dimiliki | Stok Terjual | Stok Masuk | Harga Jual | Harga Beli | Harga Pasar | Tokopedia Keyword | Keterangan |

## Testing

1. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

2. Aplikasi akan otomatis membuat header di sheet jika belum ada

3. Test CRUD operations melalui web interface

## Troubleshooting

### Error: "Failed to initialize Google Sheets service"

- Pastikan file `credentials.json` ada dan valid
- Pastikan Google Sheets API sudah diaktifkan
- Pastikan service account sudah di-share ke sheet

### Error: "Permission denied"

- Pastikan service account email sudah di-share ke sheet dengan permission "Editor"
- Pastikan Sheet ID benar

### Error: "Sheet not found"

- Pastikan Sheet ID benar
- Pastikan sheet sudah dibuat dan bisa diakses

## Mock Mode

Jika Google Sheets API tidak bisa diakses, aplikasi akan otomatis menggunakan mock service untuk testing.
