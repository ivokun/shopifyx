-- Drop the trigger
DROP TRIGGER IF EXISTS product_stock_history_insert_trigger ON product_stock;

-- Drop the trigger function
DROP FUNCTION IF EXISTS product_stock_history_insert_trigger();

-- Drop the stock_history table
DROP TABLE IF EXISTS product_stock_history;

-- Drop the current_stock table
DROP TABLE IF EXISTS product_stock;
