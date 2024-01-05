CREATE TABLE IF NOT EXISTS posts (
    post_id VARCHAR(255) PRIMARY KEY,
    source_id VARCHAR(255),
    source_name VARCHAR(255),
    author VARCHAR(255),
    title TEXT,
    description TEXT,
    url TEXT,
    url_to_image TEXT,
    published_at TIMESTAMP,
    content TEXT
);