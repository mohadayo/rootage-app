CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
