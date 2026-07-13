-- Add profiles table — one-to-one with users, holds identity/presentation data
CREATE TABLE IF NOT EXISTS profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL,
    headline            VARCHAR(150) NOT NULL,
    company_name        VARCHAR(150) NOT NULL,
    bio                 TEXT,
    skills              TEXT[],
    linkedin_url        VARCHAR(255),
    profile_photo_url   VARCHAR(255),
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Drop identity fields from users — they now live in profiles
ALTER TABLE users DROP COLUMN IF EXISTS name;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS skills;