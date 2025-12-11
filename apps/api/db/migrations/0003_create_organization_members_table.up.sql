CREATE TABLE organization_members (
  id SERIAL PRIMARY KEY,
  organization_id INT REFERENCES organizations(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  role VARCHAR(255),
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX org_member_unique
ON organization_members(organization_id, user_id);
