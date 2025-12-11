CREATE TABLE board_column (
  id SERIAL PRIMARY KEY,
  board_id INT REFERENCES boards(id) ON DELETE CASCADE,
  name VARCHAR(255),
  position INT
);
