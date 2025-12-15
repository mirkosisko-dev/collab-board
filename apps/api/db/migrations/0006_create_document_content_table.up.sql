CREATE TABLE document_content (
  document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
  ydoc_state BYTEA,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);
