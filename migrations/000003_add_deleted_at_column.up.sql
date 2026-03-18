ALTER TABLE chat_members ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX idx_members_deleted_at ON chat_members (deleted_at);