CREATE TABLE IF NOT EXISTS transactions (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  user_id INTEGER NOT NULL,
  code VARCHAR(255) NOT NULL,
  amount INTEGER NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
  -- FOREIGN KEY (user_id) REFERENCES users (id),
  -- UNIQUE (user_id, created_at),
  -- CHECK (amount > 0),
  -- CHECK (created_at > '1970-01-01 00:00:00')
);