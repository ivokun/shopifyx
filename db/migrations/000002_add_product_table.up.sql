BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_conditions') THEN
        CREATE TYPE valid_conditions AS ENUM ('new', 'used');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS products(
  id TEXT PRIMARY KEY,
  name VARCHAR (60) NOT NULL,
  imageUrl TEXT NOT NULL,
  condition valid_conditions NOT NULL,
  is_purchaseable BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),
  deleted_at TIMESTAMP NULL
);


CREATE TABLE IF NOT EXISTS product_tags (
  product_id TEXT UNIQUE NOT NULL REFERENCES products(id),
  tags TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX product_tags_id_tags ON product_tags USING gin (tags);


COMMIT;
