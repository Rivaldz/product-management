-- Seed companies
INSERT INTO companies (id, name, status) VALUES 
('00000000-0000-0000-0000-000000000001', 'PT Maju Jaya', 'ACTIVE'),
('00000000-0000-0000-0000-000000000999', 'PT Tidak Dimiliki', 'ACTIVE')
ON CONFLICT (id) DO UPDATE SET 
  name = EXCLUDED.name, 
  status = EXCLUDED.status;

-- Seed items
INSERT INTO items (id, company_id, code, name, type, price, category_name, status) VALUES 
('00000000-0000-0000-0000-000000001001', '00000000-0000-0000-0000-000000000001', 'ITEM-001', 'Paket Konsultasi Pajak', 'SERVICE', 1500000.0000, 'Konsultasi', 'ACTIVE'),
('00000000-0000-0000-0000-000000001002', '00000000-0000-0000-0000-000000000001', 'ITEM-002', 'Software Akuntansi Basic', 'PRODUCT', 2500000.0000, 'Software', 'ACTIVE'),
('00000000-0000-0000-0000-000000001003', '00000000-0000-0000-0000-000000000001', 'ITEM-003', 'Layanan Training Pajak', 'SERVICE', 1000000.0000, 'Training', 'ARCHIVED'),
('00000000-0000-0000-0000-000000001999', '00000000-0000-0000-0000-000000000999', 'ITEM-001', 'Produk Company Lain', 'PRODUCT', 500000.0000, 'Umum', 'ACTIVE')
ON CONFLICT (id) DO UPDATE SET 
  company_id = EXCLUDED.company_id, 
  code = EXCLUDED.code, 
  name = EXCLUDED.name, 
  type = EXCLUDED.type, 
  price = EXCLUDED.price, 
  category_name = EXCLUDED.category_name, 
  status = EXCLUDED.status;
