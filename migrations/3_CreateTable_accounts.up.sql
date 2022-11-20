CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id uuid NOT NULL REFERENCES clients(id) ,
  currency_id INTEGER NOT NULL REFERENCES currencies(id),
  balance DECIMAL(10, 4) DEFAULT 0
);