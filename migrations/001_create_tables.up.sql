-- Create novels table
CREATE TABLE IF NOT EXISTS novels (
    ncode VARCHAR(20) PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    author_url TEXT,
    description TEXT,
    novel_url TEXT NOT NULL,
    episode_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create episodes table
CREATE TABLE IF NOT EXISTS episodes (
    ncode VARCHAR(20) NOT NULL,
    episode_number INT NOT NULL,
    subtitle TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT,
    posted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ncode, episode_number),
    FOREIGN KEY (ncode) REFERENCES novels(ncode) ON DELETE CASCADE
);
