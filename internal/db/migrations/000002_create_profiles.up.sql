-- One-to-one with users: identity/presentation data lives here,
-- auth data stays on users.
CREATE TABLE profiles (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name               VARCHAR(100) NOT NULL,
    headline           VARCHAR(150) NOT NULL,
    company_name       VARCHAR(150) NOT NULL,
    bio                TEXT,
    skills             TEXT[],
    linkedin_url       VARCHAR(255),
    profile_photo_url  VARCHAR(255),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);