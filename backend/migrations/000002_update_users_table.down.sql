ALTER TABLE users RENAME COLUMN first_name TO name;
ALTER TABLE users RENAME COLUMN is_email_verified TO is_verified;
ALTER TABLE users 
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS date_of_birth,
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS profile_image_url,
    DROP COLUMN IF EXISTS rating_average,
    DROP COLUMN IF EXISTS rating_count,
    DROP COLUMN IF EXISTS is_phone_verified,
    DROP COLUMN IF EXISTS is_id_verified,
    DROP COLUMN IF EXISTS email_verification_token,
    DROP COLUMN IF EXISTS email_verification_expiry,
    DROP COLUMN IF EXISTS password_reset_token,
    DROP COLUMN IF EXISTS password_reset_token_expiry,
    DROP COLUMN IF EXISTS last_seen_at,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS updated_at;