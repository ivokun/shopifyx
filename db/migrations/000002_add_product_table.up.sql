BEGIN;

-- Add valid_conditions enums
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_conditions') THEN
        CREATE TYPE valid_conditions AS ENUM ('new', 'used');
    END IF;
END
$$;

-- Add products table
CREATE TABLE IF NOT EXISTS products(
  id TEXT PRIMARY KEY,
  name VARCHAR (60) NOT NULL,
  imageUrl TEXT NOT NULL,
  condition valid_conditions NOT NULL,
  is_purchaseable BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),
  deleted_at TIMESTAMP NULL,

  seller_id TEXT NOT NULL REFERENCES users(id)
);

-- Add product_tags table
CREATE TABLE IF NOT EXISTS product_tags (
  product_id TEXT UNIQUE NOT NULL REFERENCES products(id),
  tags TEXT[] NOT NULL DEFAULT '{}'
);

-- Add index on tags
CREATE INDEX product_tags_id_tags ON product_tags USING GIN (tags);


COMMIT;
