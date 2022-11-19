CREATE TABLE IF NOT EXISTS transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  type_id SMALLINT NOT NULL,
  source_id uuid REFERENCES accounts(id) ON DELETE CASCADE,
  target_id uuid REFERENCES accounts(id) ON DELETE CASCADE,
  amount DECIMAL(10, 4) NOT NULL
);