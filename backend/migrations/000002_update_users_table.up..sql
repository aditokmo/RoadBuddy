ALTER TABLE users RENAME COLUMN name TO first_name;
ALTER TABLE users 
    ADD COLUMN last_name TEXT NOT NULL DEFAULT '',
    ADD COLUMN date_of_birth TIMESTAMPTZ,
    ADD COLUMN phone_number TEXT,
    ADD COLUMN profile_image_url TEXT,
    ADD COLUMN rating_average FLOAT8 NOT NULL DEFAULT 0.0,
    ADD COLUMN rating_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_id_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN email_verification_token TEXT,
    ADD COLUMN email_verification_expiry TIMESTAMPTZ,
    ADD COLUMN password_reset_token TEXT,
    ADD COLUMN password_reset_token_expiry TIMESTAMPTZ,
    ADD COLUMN last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE users RENAME COLUMN is_verified TO is_email_verified;