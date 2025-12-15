CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE(email);
ALTER TABLE users ADD CONSTRAINT users_email_lowercase CHECK (email = LOWER(email));
