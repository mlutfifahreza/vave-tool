DROP INDEX IF EXISTS idx_subcategories_created_at;
DROP INDEX IF EXISTS idx_subcategories_is_active;
DROP INDEX IF EXISTS idx_subcategories_name;
DROP INDEX IF EXISTS idx_subcategories_category_id;
ALTER TABLE subcategories DROP CONSTRAINT IF EXISTS fk_subcategories_category_id;
DROP TABLE IF EXISTS subcategories;