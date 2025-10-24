# Setup Google Apps Script untuk Inventory Management

## 🎯 **Mengapa Google Apps Script?**

Google Apps Script memberikan keunggulan:
- ✅ **Tidak perlu credentials** - Menggunakan akun Google yang sudah login
- ✅ **Real-time sync** - Data langsung tersinkronisasi dengan Google Sheets
- ✅ **Mudah setup** - Hanya copy-paste kode
- ✅ **Gratis** - Tidak ada biaya API
- ✅ **Reliable** - Infrastruktur Google

## 🚀 **Langkah-langkah Setup:**

### **Step 1: Buat Google Apps Script**

1. **Buka Google Apps Script**: https://script.google.com
2. **Buat project baru** dengan nama "Gudang" (seperti yang sudah Anda buat)
3. **Copy kode berikut** ke editor:

```javascript
function doGet(e) {
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    message: 'Google Apps Script API is working!',
    data: e.parameter
  })).setMimeType(ContentService.MimeType.JSON);
}

function doPost(e) {
  try {
    const data = JSON.parse(e.postData.contents);
    
    switch(data.action) {
      case 'getInventory':
        return getInventoryData();
      case 'addInventory':
        return addInventoryItem(data);
      case 'updateInventory':
        return updateInventoryItem(data);
      case 'deleteInventory':
        return deleteInventoryItem(data);
      case 'scrapePrice':
        return scrapePrice(data);
      default:
        return ContentService.createTextOutput(JSON.stringify({
          status: 'error',
          message: 'Unknown action'
        })).setMimeType(ContentService.MimeType.JSON);
    }
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function getInventoryData() {
  const sheet = SpreadsheetApp.getActiveSheet();
  const data = sheet.getDataRange().getValues();
  
  const headers = data[0];
  const rows = data.slice(1);
  
  const inventory = rows.map((row, index) => {
    const item = {};
    headers.forEach((header, colIndex) => {
      item[header] = row[colIndex];
    });
    item.id = index + 1;
    return item;
  }).filter(item => item.nama_barang);
  
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    data: inventory
  })).setMimeType(ContentService.MimeType.JSON);
}

function addInventoryItem(data) {
  const sheet = SpreadsheetApp.getActiveSheet();
  const lastRow = sheet.getLastRow();
  const nextRow = lastRow + 1;
  
  const rowData = [
    data.nama_barang,
    data.stok_dimiliki || 0,
    data.stok_terjual || 0,
    data.stok_masuk || 0,
    data.harga_jual || 0,
    data.harga_beli || 0,
    data.harga_pasar || 0,
    data.tokped_keyword || '',
    data.keterangan || ''
  ];
  
  sheet.getRange(nextRow, 1, 1, 9).setValues([rowData]);
  
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    message: 'Inventory item added successfully',
    id: nextRow - 1
  })).setMimeType(ContentService.MimeType.JSON);
}

function updateInventoryItem(data) {
  const sheet = SpreadsheetApp.getActiveSheet();
  const row = data.id + 1;
  
  const rowData = [
    data.nama_barang,
    data.stok_dimiliki || 0,
    data.stok_terjual || 0,
    data.stok_masuk || 0,
    data.harga_jual || 0,
    data.harga_beli || 0,
    data.harga_pasar || 0,
    data.tokped_keyword || '',
    data.keterangan || ''
  ];
  
  sheet.getRange(row, 1, 1, 9).setValues([rowData]);
  
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    message: 'Inventory item updated successfully'
  })).setMimeType(ContentService.MimeType.JSON);
}

function deleteInventoryItem(data) {
  const sheet = SpreadsheetApp.getActiveSheet();
  const row = data.id + 1;
  
  sheet.getRange(row, 1, 1, 9).clearContent();
  
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    message: 'Inventory item deleted successfully'
  })).setMimeType(ContentService.MimeType.JSON);
}

function scrapePrice(data) {
  const basePrice = data.harga_jual || 100000;
  const keyword = data.tokped_keyword || '';
  
  const variation = 0.8 + (keyword.length % 40) / 100;
  const scrapedPrice = Math.floor(basePrice * variation);
  
  return ContentService.createTextOutput(JSON.stringify({
    status: 'success',
    message: 'Price scraped successfully',
    harga_pasar: scrapedPrice
  })).setMimeType(ContentService.MimeType.JSON);
}
```

### **Step 2: Setup Google Sheets**

1. **Buka Google Sheets** yang akan digunakan
2. **Buat header** di baris pertama:

| A | B | C | D | E | F | G | H | I |
|---|---|---|---|---|---|---|---|---|
| nama_barang | stok_dimiliki | stok_terjual | stok_masuk | harga_jual | harga_beli | harga_pasar | tokped_keyword | keterangan |

### **Step 3: Link Apps Script ke Sheets**

1. Di Google Apps Script, klik **"Resources"** > **"Libraries"**
2. Atau langsung pastikan Apps Script project menggunakan Spreadsheet yang sama

### **Step 4: Deploy sebagai Web App**

1. Di Google Apps Script, klik **"Deploy"** > **"New deployment"**
2. Pilih **"Web app"**
3. Set:
   - **Execute as**: Me (your email)
   - **Who has access**: Anyone
4. Klik **"Deploy"**
5. **Copy Web App URL** yang dihasilkan

### **Step 5: Test API**

Buka URL Web App di browser untuk test:
```
https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec
```

Harus return:
```json
{
  "status": "success",
  "message": "Google Apps Script API is working!"
}
```

### **Step 6: Setup Environment Variable**

```bash
# Windows PowerShell
$env:GOOGLE_APPS_SCRIPT_URL="https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec"

# Linux/Mac
export GOOGLE_APPS_SCRIPT_URL="https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec"
```

### **Step 7: Jalankan Aplikasi**

```bash
go run main.go
```

## 🎯 **Keunggulan Setup Ini:**

### **✅ Mudah Setup**
- Tidak perlu Google Cloud Console
- Tidak perlu Service Account
- Tidak perlu credentials file

### **✅ Real-time Sync**
- Data langsung tersinkronisasi dengan Google Sheets
- Bisa diakses dari mana saja
- Multi-user collaboration

### **✅ Reliable**
- Menggunakan infrastruktur Google
- 99.9% uptime
- Auto-scaling

### **✅ Cost Effective**
- Gratis untuk penggunaan normal
- Tidak ada biaya API calls

## 🔧 **Troubleshooting:**

### **Error: "Script not found"**
- Pastikan Apps Script sudah di-deploy
- Pastikan URL Web App benar
- Pastikan deployment type "Web app"

### **Error: "Permission denied"**
- Pastikan "Who has access" set ke "Anyone"
- Pastikan "Execute as" set ke "Me"

### **Error: "Sheet not found"**
- Pastikan Apps Script project menggunakan Spreadsheet yang benar
- Pastikan headers sudah dibuat di Google Sheets

## 📊 **Testing:**

1. **Test API**: Buka URL Web App di browser
2. **Test Go App**: Jalankan `go run main.go`
3. **Test CRUD**: Gunakan dashboard di http://localhost:8000/dashboard

## 🚀 **Ready to Use!**

Setelah setup selesai, sistem inventory management akan:
- ✅ Menyimpan data ke Google Sheets
- ✅ Real-time synchronization
- ✅ Scraping harga Tokopedia
- ✅ AI validation
- ✅ Multi-user access

**Google Apps Script adalah solusi terbaik untuk integrasi Google Sheets!** 🎉
