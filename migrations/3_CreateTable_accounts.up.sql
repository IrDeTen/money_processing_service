CREATE TABLE IF NOT EXISTS accounts (
  id uuid PRIMARY KEY,
  client_id uuid NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
  currency_id INTEGER NOT NULL REFERENCES currencys(id) ON DELETE CASCADE,
  ballance DECIMAL(10, 4) NOT NULL
);