-- PostGIS gives us the GEOGRAPHY type and spatial functions (ST_DWithin, ST_Distance).
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR(255) NOT NULL UNIQUE,   -- UNIQUE already creates a B-tree index on email
    password    VARCHAR(255) NOT NULL,          -- bcrypt hash, never plaintext
    location    GEOGRAPHY(POINT, 4326),         -- NULL until the user calls PUT /location
    is_visible  BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- GIST = spatial index. Lets ST_DWithin skip every region that can't
-- overlap the search circle instead of checking every row.
CREATE INDEX idx_users_location ON users USING GIST (location);