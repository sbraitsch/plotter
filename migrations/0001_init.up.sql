CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE player_mappings (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    from_num INT NOT NULL CHECK (from_num BETWEEN 1 AND 53),
    to_num   INT NOT NULL CHECK (to_num BETWEEN 1 AND 53),
    UNIQUE (player_id, from_num),
    UNIQUE (player_id, to_num)
);
