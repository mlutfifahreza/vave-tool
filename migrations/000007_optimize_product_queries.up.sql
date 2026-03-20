CREATE INDEX CONCURRENTLY idx_products_active_category_subcategory_created 
ON products (is_active, category_id, subcategory_id, created_at DESC)
WHERE is_active = true;

CREATE INDEX CONCURRENTLY idx_products_active_category_price_created 
ON products (is_active, category_id, price, created_at DESC)
WHERE is_active = true;

CREATE INDEX CONCURRENTLY idx_products_active_price_created 
ON products (is_active, price, created_at DESC)
WHERE is_active = true;
