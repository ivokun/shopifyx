CREATE TABLE IF NOT EXISTS users(
  id TEXT PRIMARY KEY,
  username VARCHAR (50) UNIQUE NOT NULL,
  password VARCHAR (50) NOT NULL,
  created_at timestamp default (timezone('utc', now())),
  deleted_at timestamp NULL,
);
