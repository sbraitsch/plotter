-- Add a boolean column to indicate if a player is an admin
ALTER TABLE players
ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

