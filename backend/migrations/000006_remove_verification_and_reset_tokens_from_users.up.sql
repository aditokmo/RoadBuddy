ALTER TABLE users 
    DROP COLUMN IF EXISTS email_verification_token,
    DROP COLUMN IF EXISTS email_verification_expiry,
    DROP COLUMN IF EXISTS password_reset_token,
    DROP COLUMN IF EXISTS password_reset_token_expiry;