-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('talent', 'client')),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS talent_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    talent_type TEXT NOT NULL DEFAULT 'model',
    display_name TEXT NOT NULL,
    city TEXT NOT NULL,
    age INT NOT NULL CHECK (age >= 18),
    height_cm INT,
    bio TEXT NOT NULL DEFAULT '',
    rate_min INT NOT NULL DEFAULT 0,
    rate_max INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_talent_profiles_city_type ON talent_profiles (city, talent_type);

CREATE TABLE IF NOT EXISTS client_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT NOT NULL,
    org_name TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS profile_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    talent_id UUID NOT NULL REFERENCES talent_profiles(user_id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_profile_photos_talent_pos ON profile_photos (talent_id, position);

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    talent_type TEXT NOT NULL DEFAULT 'model',
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL,
    budget BIGINT NOT NULL DEFAULT 0,
    shoot_date DATE NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('open', 'filled', 'closed')) DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_jobs_feed ON jobs (city, talent_type, status, shoot_date);

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    talent_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('applied', 'shortlisted', 'rejected', 'withdrawn')) DEFAULT 'applied',
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_applications_job_talent UNIQUE (job_id, talent_id)
);

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_matches_application UNIQUE (application_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_messages_match_sent ON messages (match_id, sent_at);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    amount_paise BIGINT NOT NULL CHECK (amount_paise > 0),
    status TEXT NOT NULL CHECK (status IN ('pending', 'held', 'released', 'refunded', 'failed')) DEFAULT 'pending',
    gateway_ref TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone TEXT NOT NULL,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_otps_phone ON otps (phone);

-- +goose Down
DROP TABLE IF EXISTS otps;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS profile_photos;
DROP TABLE IF EXISTS client_profiles;
DROP TABLE IF EXISTS talent_profiles;
DROP TABLE IF EXISTS users;

