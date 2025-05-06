CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  account_id TEXT,
  balance NUMERIC,
  nonce BIGINT,
  code_hash TEXT,
  storage_root TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

