CREATE TABLE IF NOT EXISTS members(
  id SERIAL PRIMARY KEY NOT NULL,
  slack_user_id VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cards(
  id SERIAL PRIMARY KEY NOT NULL,
  sender_member_id INT REFERENCES members(id) NOT NULL,
  distination_member_id INT REFERENCES members(id) NOT NULL,
  point INT NOT NULL,
  message VARCHAR(400) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
