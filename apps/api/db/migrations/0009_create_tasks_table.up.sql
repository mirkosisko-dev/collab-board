CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  board_id INT REFERENCES boards(id) ON DELETE CASCADE,
  column_id INT REFERENCES board_column(id) ON DELETE CASCADE,
  assignee_id INT REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255),
  description VARCHAR(255),
  position INT,
  created_by INT REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
