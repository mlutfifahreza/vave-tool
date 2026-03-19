DROP INDEX IF EXISTS idx_products_subcategory_id;
DROP INDEX IF EXISTS idx_products_category_id;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_subcategory_id;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_category_id;
ALTER TABLE products DROP COLUMN subcategory_id;
ALTER TABLE products DROP COLUMN category_id;
ALTER TABLE products ADD COLUMN category VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);