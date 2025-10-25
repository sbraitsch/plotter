CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS communities (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE,
    realm VARCHAR(50),
    officer_rank INT DEFAULT 0,
    member_rank INT DEFAULT 1,
    locked BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    battletag VARCHAR(50) PRIMARY KEY,
    char VARCHAR(13),
    community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    community_rank INT DEFAULT 100,
    session_id UUID UNIQUE,
    access_token TEXT NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE plot_mappings (
    id SERIAL PRIMARY KEY,
    battletag VARCHAR(50) NOT NULL REFERENCES users(battletag) ON DELETE CASCADE,
    plot_id INT NOT NULL CHECK (plot_id BETWEEN 1 AND 53),
    priority   INT NOT NULL CHECK (priority BETWEEN 1 AND 53),
    UNIQUE (battletag, plot_id),
    UNIQUE (battletag, priority)
);

CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    battletag VARCHAR(50) NOT NULL REFERENCES users(battletag) ON DELETE CASCADE,
    char VARCHAR(50),
    community_id UUID REFERENCES communities(id) ON DELETE CASCADE,
    plot_id INT NOT NULL,
    plot_score INT NOT NULL,
    UNIQUE (battletag)
);
