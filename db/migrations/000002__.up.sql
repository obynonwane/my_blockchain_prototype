CREATE TABLE accounts (
  account_id TEXT PRIMARY KEY,
  balance NUMERIC,
  nonce BIGINT,
  code_hash TEXT,
  storage_root TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);