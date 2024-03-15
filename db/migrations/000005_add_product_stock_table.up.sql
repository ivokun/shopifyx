-- Add product_stock table
CREATE TABLE product_stock(
  id TEXT PRIMARY KEY,
  quantity INT NOT NULL,
  created_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),
  modified_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),

  product_id TEXT NOT NULL REFERENCES products(id)
);

-- Add product_stock_history
CREATE TABLE product_stock_history(
  id SERIAL PRIMARY KEY,
  stock_level INTEGER NOT NULL,
  change_date TIMESTAMP NOT NULL,
  change_type VARCHAR(10) NOT NULL,

  product_id TEXT NOT NULL REFERENCES products(id)
);

-- Add stock_history_insert_trigger function
CREATE OR REPLACE FUNCTION product_stock_history_insert_trigger()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO product_stock_history (product_id, stock_level, change_date, change_type)
    VALUES (NEW.product_id, NEW.quantity, NOW(), 'update');
    RETURN NEW;
END;
$$
 LANGUAGE plpgsql;

-- Add trigger to update product_stock_history based on insert events on product_stock
CREATE TRIGGER product_stock_history_insert_trigger
AFTER UPDATE OF quantity ON product_stock
FOR EACH ROW
EXECUTE FUNCTION product_stock_history_insert_trigger();


