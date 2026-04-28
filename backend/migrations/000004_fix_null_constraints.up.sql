UPDATE users SET phone_number = '' WHERE phone_number IS NULL;
UPDATE users SET profile_image_url = '' WHERE profile_image_url IS NULL;

ALTER TABLE users ALTER COLUMN phone_number SET NOT NULL;
ALTER TABLE users ALTER COLUMN phone_number SET DEFAULT '';

ALTER TABLE users ALTER COLUMN profile_image_url SET NOT NULL;
ALTER TABLE users ALTER COLUMN profile_image_url SET DEFAULT '';