CREATE SEQUENCE IF NOT EXISTS products_code_seq;

CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    code  TEXT NOT NULL DEFAULT (nextval('products_code_seq'))::text,
    name  TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL DEFAULT 0,
    CONSTRAINT products_code_key UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS products_name_idx ON products (name);
