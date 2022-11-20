CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id uuid NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
  currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
  ballance DECIMAL(10, 4) DEFAULT 0
);