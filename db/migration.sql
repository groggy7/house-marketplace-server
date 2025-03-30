CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
)

CREATE TABLE messages (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    chat_id TEXT GENERATED ALWAYS AS (
    LEAST(sender_id, receiver_id) || '-' || GREATEST(sender_id, receiver_id)
    ) STORED,
    message TEXT NOT NULL,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    read_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_chat ON messages (chat_id, created_at);