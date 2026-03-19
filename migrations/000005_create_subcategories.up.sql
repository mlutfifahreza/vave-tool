CREATE TABLE IF NOT EXISTS subcategories (
    id VARCHAR(100) PRIMARY KEY,
    category_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, name)
);

ALTER TABLE subcategories ADD CONSTRAINT fk_subcategories_category_id FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_subcategories_category_id ON subcategories(category_id);
CREATE INDEX IF NOT EXISTS idx_subcategories_name ON subcategories(name);
CREATE INDEX IF NOT EXISTS idx_subcategories_is_active ON subcategories(is_active);
CREATE INDEX IF NOT EXISTS idx_subcategories_created_at ON subcategories(created_at);