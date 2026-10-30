-- Local test data. NEVER run against production.
-- Every user's password is: password123
-- Run AFTER the server has started once (so migrations have created the tables):
--   psql -d proximate_db -f scripts/seed.sql

BEGIN;

TRUNCATE users CASCADE;  -- wipes users + their profiles + refresh tokens

INSERT INTO users (email, password, location, is_visible) VALUES
  -- Bob and Charlie Lee share an IDENTICAL point on purpose: equal distance
  -- from any searcher, so they exercise the (distance, id) cursor tiebreak.
  ('bob123@test.com',     '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5950, 12.9720), 4326)::geography, true),
  ('charlie123@test.com', '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5950, 12.9720), 4326)::geography, true),
  ('hetul123@test.com',   '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5946, 12.9716), 4326)::geography, true),
  ('evan@test.com',       '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5955, 12.9725), 4326)::geography, true),
  ('dana@test.com',       '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5940, 12.9710), 4326)::geography, true),
  -- Should NEVER appear in /nearby, each for a different reason:
  ('charlie@test.com',    '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5948, 12.9718), 4326)::geography, true),   -- no profile
  ('ghost@test.com',      '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(77.5947, 12.9717), 4326)::geography, false),  -- is_visible = false
  ('raj@test.com',        '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', ST_SetSRID(ST_MakePoint(72.5812, 23.0328), 4326)::geography, true),   -- in Ahmedabad, far outside radius
  ('nolocation@test.com', '$2a$10$WiSjjL.1HMzv/JhNjpXIMO.jGHxlPgxNelxOTDAeVE2APh7qrkPAu', NULL, true);                                                       -- never set a location

INSERT INTO profiles (user_id, name, headline, company_name, skills)
SELECT u.id, p.name, p.headline, p.company, p.skills
FROM (VALUES
  ('bob123@test.com',     'Bob Smith',   'Product Manager',          'Acme Corp',     ARRAY['roadmaps','sql']),
  ('charlie123@test.com', 'Charlie Lee', 'Backend Engineer',         'Initech',       ARRAY['go','postgres','redis']),
  ('hetul123@test.com',   'Hetul Mehta', 'Senior Backend Developer', 'Proximate Inc', ARRAY['go','postgis','system design']),
  ('evan@test.com',       'Evan Ross',   'DevOps Lead',              'Umbrella',      ARRAY['docker','kubernetes']),
  ('dana@test.com',       'Dana Kapoor', 'Data Scientist',           'Globex',        ARRAY['python','ml']),
  ('ghost@test.com',      'Ghost User',  'Hidden',                   'Nowhere',       ARRAY['stealth']),
  ('raj@test.com',        'Raj Patel',   'Frontend Engineer',        'Ahmedabad Labs',ARRAY['javascript','css']),
  ('nolocation@test.com', 'Nora Null',   'Designer',                 'Nomad Studio',  ARRAY['figma'])
) AS p(email, name, headline, company, skills)
JOIN users u ON u.email = p.email;

COMMIT;