CREATE TABLE organization_invite (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  invited_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  invited_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role organization_role NOT NULL DEFAULT 'owner',
  status organization_invite_status NOT NULL DEFAULT 'pending', 
  expires_at TIMESTAMP WITHOUT TIME ZONE,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  responded_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX organization_invite_unique_pending
ON organization_invite (organization_id, invited_user_id)
WHERE status = 'pending';
