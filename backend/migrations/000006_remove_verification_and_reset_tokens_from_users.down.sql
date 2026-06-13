ALTER TABLE users 
    ADD COLUMN email_verification_token TEXT,
    ADD COLUMN email_verification_expiry TIMESTAMPTZ,
    ADD COLUMN password_reset_token TEXT,
    ADD COLUMN password_reset_token_expiry TIMESTAMPTZ;