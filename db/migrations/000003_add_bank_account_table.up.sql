CREATE TABLE IF NOT EXISTS bank_accounts(
  id TEXT PRIMARY KEY,
  bank_name TEXT NOT NULL,
  bank_account_name INT NOT NULL,
  bank_account_number INT NOT NULL,
  created_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),
  modified_at TIMESTAMP DEFAULT (TIMEZONE('UTC', NOW())),
  deleted_at TIMESTAMP,

  user_id TEXT NOT NULL REFERENCES users(id)
);
