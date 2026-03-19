DROP INDEX IF EXISTS idx_products_updated_by;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_updated_by;
ALTER TABLE products DROP COLUMN IF EXISTS updated_by;
