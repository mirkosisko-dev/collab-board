CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  board_id INT REFERENCES boards(id) ON DELETE CASCADE,
  organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  document_id INT REFERENCES documents(id) ON DELETE CASCADE,
  content VARCHAR(255),
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
