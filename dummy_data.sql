-- Enable pgcrypto if not yet active
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Clear old data for re-runs
TRUNCATE assignments, plot_mappings, users, communities RESTART IDENTITY CASCADE;

-- 1️⃣ Insert one community
INSERT INTO communities (name, realm, officer_rank, member_rank, locked)
VALUES ('Test Community', 'NA-1', 3, 3, FALSE)
RETURNING id INTO TEMP TABLE tmp_community;

-- 2️⃣ Generate 50 users in that community
DO $$
DECLARE
    comm_id UUID;
    i INT;
    battletag TEXT;
BEGIN
    SELECT id INTO comm_id FROM tmp_community;

    FOR i IN 1..50 LOOP
        battletag := format('Player#%s', lpad(i::text, 4, '0'));
        INSERT INTO users (
            battletag, community_id, community_rank, session_id,
            access_token, expiry, created_at, updated_at
        )
        VALUES (
            battletag,
            comm_id,
            (RANDOM() * 4)::INT,
            gen_random_uuid(),
            encode(gen_random_bytes(16), 'hex'),
            NOW() + INTERVAL '1 day',
            NOW(),
            NOW()
        );
    END LOOP;
END $$;

-- 3️⃣ Insert 5–10 plot mappings per user
DO $$
DECLARE
    u RECORD;
    num_mappings INT;
    plots INT[];
    priorities INT[];
    j INT;
BEGIN
    FOR u IN SELECT battletag FROM users LOOP
        -- Random number of mappings between 5 and 10
        num_mappings := 5 + (random() * 5)::INT;

        -- Random plots between 1 and 53, with overlap (many users share plots)
        plots := ARRAY(
            SELECT (1 + floor(random() * 53))::INT
            FROM generate_series(1, num_mappings)
        );

        -- Random priority 1–10 (so overlaps happen easily)
        priorities := ARRAY(
            SELECT (1 + floor(random() * 10))::INT
            FROM generate_series(1, num_mappings)
        );

        FOR j IN 1..num_mappings LOOP
            BEGIN
                INSERT INTO plot_mappings (battletag, plot_id, priority)
                VALUES (u.battletag, plots[j], priorities[j]);
            EXCEPTION WHEN unique_violation THEN
                -- Skip duplicates
                CONTINUE;
            END;
        END LOOP;
    END LOOP;
END $$;
