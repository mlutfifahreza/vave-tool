ALTER TABLE products ADD COLUMN updated_by UUID;
ALTER TABLE products ADD CONSTRAINT fk_products_updated_by FOREIGN KEY (updated_by) REFERENCES clients(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_products_updated_by ON products(updated_by);
