CREATE TABLE IF NOT EXISTS payments(
  id TEXT PRIMARY KEY,
  payment_proof_image_url TEXT NOT NULL,
  quantity INT NOT NULL,
  created_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),

  product_id TEXT NOT NULL REFERENCES products(id),
  bank_account_id TEXT NOT NULL REFERENCES bank_accounts(id)
);
