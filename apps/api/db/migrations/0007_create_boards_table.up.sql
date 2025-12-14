CREATE TABLE boards (
  id SERIAL PRIMARY KEY,
  organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
  name VARCHAR(255),
  created_by INT REFERENCES users(id) ON DELETE CASCADE
);
