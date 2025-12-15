CREATE TABLE organization_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role organization_role NOT NULL DEFAULT 'owner',
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

ALTER TABLE organization_members
ADD CONSTRAINT organization_member_org_user_unique UNIQUE (organization_id, user_id);

CREATE UNIQUE INDEX org_member_unique
ON organization_members(organization_id, user_id);
