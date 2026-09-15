CREATE TYPE item_type AS ENUM ('story', 'comment', 'job', 'poll', 'pollopt');

CREATE TABLE IF NOT EXISTS items (
    id BIGINT PRIMARY KEY,
    deleted BOOLEAN DEFAULT FALSE,
    type item_type NOT NULL,
    by_user VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    text TEXT,
    dead BOOLEAN DEFAULT FALSE,
    parent_id BIGINT REFERENCES items(id),
    poll_id BIGINT REFERENCES items(id),
    url TEXT,
    score INT DEFAULT 0,
    title VARCHAR(500),
    descendants INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);