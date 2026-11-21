-- 001: products table (aligned with models.Product)

CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS products_name_idx ON products (name);
