-- Enable PostGIS (already done, but idempotent so safe to re-run)
CREATE EXTENSION IF NOT EXISTS postgis;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    bio         TEXT,
    skills      TEXT[],
    location    GEOGRAPHY(POINT, 4326),
    is_visible  BOOLEAN DEFAULT true,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for fast geospatial queries
CREATE INDEX IF NOT EXISTS idx_users_location ON users USING GIST(location);

-- Index for email lookups (login)
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- Redundant because whenever we use UNIQUE keyword, automatically B tree index is created in PostgreSQL