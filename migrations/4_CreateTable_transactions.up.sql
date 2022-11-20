CREATE TABLE IF NOT EXISTS transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  creation_date TIMESTAMP NOT NULL,
  type_id SMALLINT NOT NULL,
  source_id uuid REFERENCES accounts(id),
  target_id uuid REFERENCES accounts(id),
  amount DECIMAL(10, 4) NOT NULL
);
