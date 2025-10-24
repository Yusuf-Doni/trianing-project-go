// Google Apps Script untuk Inventory Management
// Pastikan script ini terhubung dengan Spreadsheet yang benar

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
          message: 'Unknown action: ' + data.action
        })).setMimeType(ContentService.MimeType.JSON);
    }
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'Error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function getInventoryData() {
  try {
    // Pastikan menggunakan Spreadsheet yang benar
    const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
    const sheet = spreadsheet.getActiveSheet();
    
    // Debug: Log sheet info
    console.log('Sheet name: ' + sheet.getName());
    console.log('Spreadsheet name: ' + spreadsheet.getName());
    
    const data = sheet.getDataRange().getValues();
    
    if (data.length <= 1) {
      return ContentService.createTextOutput(JSON.stringify({
        status: 'success',
        data: []
      })).setMimeType(ContentService.MimeType.JSON);
    }
    
    const headers = data[0];
    const rows = data.slice(1);
    
    console.log('Headers:', headers);
    console.log('Number of rows:', rows.length);
    
    const inventory = rows.map((row, index) => {
      const item = {};
      headers.forEach((header, colIndex) => {
        item[header] = row[colIndex];
      });
      item.id = index + 1;
      return item;
    }).filter(item => item.nama_barang && item.nama_barang.toString().trim() !== '');
    
    console.log('Filtered inventory:', inventory.length);
    
    return ContentService.createTextOutput(JSON.stringify({
      status: 'success',
      data: inventory
    })).setMimeType(ContentService.MimeType.JSON);
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'getInventoryData error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function addInventoryItem(data) {
  try {
    const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
    const sheet = spreadsheet.getActiveSheet();
    
    console.log('Adding inventory item:', data);
    
    // Get the next empty row
    const lastRow = sheet.getLastRow();
    const nextRow = lastRow + 1;
    
    // Prepare data row sesuai dengan struktur Google Sheets
    const rowData = [
      data.nama_barang || '',
      data.stok_dimiliki || 0,
      data.stok_terjual || 0,
      data.stok_masuk || 0,
      data.harga_jual || 0,
      data.harga_beli || 0,
      data.harga_pasar || 0,
      data.tokped_keyword || '',
      data.keterangan || ''
    ];
    
    console.log('Row data to insert:', rowData);
    
    // Add data to sheet
    sheet.getRange(nextRow, 1, 1, 9).setValues([rowData]);
    
    console.log('Data inserted at row:', nextRow);
    
    return ContentService.createTextOutput(JSON.stringify({
      status: 'success',
      message: 'Inventory item added successfully',
      id: nextRow - 1,
      row: nextRow
    })).setMimeType(ContentService.MimeType.JSON);
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'addInventoryItem error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function updateInventoryItem(data) {
  try {
    const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
    const sheet = spreadsheet.getActiveSheet();
    const row = data.id + 1; // +1 for header row
    
    const rowData = [
      data.nama_barang || '',
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
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'updateInventoryItem error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function deleteInventoryItem(data) {
  try {
    const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
    const sheet = spreadsheet.getActiveSheet();
    const row = data.id + 1; // +1 for header row
    
    // Clear the row
    sheet.getRange(row, 1, 1, 9).clearContent();
    
    return ContentService.createTextOutput(JSON.stringify({
      status: 'success',
      message: 'Inventory item deleted successfully'
    })).setMimeType(ContentService.MimeType.JSON);
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'deleteInventoryItem error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

function scrapePrice(data) {
  try {
    const basePrice = data.harga_jual || 100000;
    const keyword = data.tokped_keyword || '';
    
    // Simulasi price scraping
    const variation = 0.8 + (keyword.length % 40) / 100;
    const scrapedPrice = Math.floor(basePrice * variation);
    
    return ContentService.createTextOutput(JSON.stringify({
      status: 'success',
      message: 'Price scraped successfully',
      harga_pasar: scrapedPrice
    })).setMimeType(ContentService.MimeType.JSON);
  } catch (error) {
    return ContentService.createTextOutput(JSON.stringify({
      status: 'error',
      message: 'scrapePrice error: ' + error.toString()
    })).setMimeType(ContentService.MimeType.JSON);
  }
}

// Test function untuk debugging
function testConnection() {
  try {
    const spreadsheet = SpreadsheetApp.getActiveSpreadsheet();
    const sheet = spreadsheet.getActiveSheet();
    
    console.log('Spreadsheet name: ' + spreadsheet.getName());
    console.log('Sheet name: ' + sheet.getName());
    console.log('Sheet ID: ' + sheet.getSheetId());
    
    return {
      status: 'success',
      spreadsheet: spreadsheet.getName(),
      sheet: sheet.getName(),
      sheetId: sheet.getSheetId()
    };
  } catch (error) {
    return {
      status: 'error',
      message: error.toString()
    };
  }
}
