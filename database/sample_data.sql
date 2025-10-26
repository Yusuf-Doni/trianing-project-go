-- Sample data untuk testing sistem inventory
-- Sesuai dengan struktur Google Sheets yang ditampilkan

INSERT INTO books (nama_barang, stok, terjual, harga,  harga_pasar, tokped_keyword, keterangan) VALUES
('Carbon Brush Starter Carry Extra', 25, 5, 150000, 145000, 'carbon brush starter carry extra', 'FS-125'),
('Laptop ASUS ROG Strix', 8, 12, 12500000, 12300000, 'laptop asus rog strix gaming', 'G15-2023'),
('Mouse Wireless Logitech MX Master', 45, 15, 850000, 820000, 'mouse wireless logitech mx master', 'MX3'),
('Keyboard Mechanical RGB', 32, 8, 1200000, 1150000, 'keyboard mechanical rgb gaming', 'K95-RGB'),
('Monitor Gaming 27 inch', 15, 5, 4500000, 4350000, 'monitor gaming 27 inch 144hz', 'VG27AQ'),
('Speaker Bluetooth JBL', 28, 12, 650000, 630000, 'speaker bluetooth jbl charge', 'Charge-4'),
('Power Bank 20000mAh', 50, 25, 250000, 240000, 'power bank 20000mah fast charging', 'PB-20000'),
('USB Cable Type-C', 100, 40, 45000, 42000, 'usb cable type c fast charging', 'UC-C1M');

-- Update harga pasar dengan simulasi scraping
UPDATE inventory SET harga_pasar = 145000 WHERE nama_barang = 'Carbon Brush Starter Carry Extra';
UPDATE inventory SET harga_pasar = 12300000 WHERE nama_barang = 'Laptop ASUS ROG Strix';
UPDATE inventory SET harga_pasar = 820000 WHERE nama_barang = 'Mouse Wireless Logitech MX Master';
UPDATE inventory SET harga_pasar = 1150000 WHERE nama_barang = 'Keyboard Mechanical RGB';
UPDATE inventory SET harga_pasar = 4350000 WHERE nama_barang = 'Monitor Gaming 27 inch';
UPDATE inventory SET harga_pasar = 630000 WHERE nama_barang = 'Speaker Bluetooth JBL';
UPDATE inventory SET harga_pasar = 240000 WHERE nama_barang = 'Power Bank 20000mAh';
UPDATE inventory SET harga_pasar = 42000 WHERE nama_barang = 'USB Cable Type-C';
