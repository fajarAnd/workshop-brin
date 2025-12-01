-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Step 1: Add new UUID column
ALTER TABLE users ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();

-- Step 2: Update all existing rows to have UUIDs
UPDATE users SET id_new = uuid_generate_v4();

-- Step 3: Make new column NOT NULL
ALTER TABLE users ALTER COLUMN id_new SET NOT NULL;

-- Step 4: Drop the primary key constraint and old column
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users DROP COLUMN id;

-- Step 5: Rename new column to id
ALTER TABLE users RENAME COLUMN id_new TO id;

-- Step 6: Add primary key constraint on UUID column
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- Step 7: Ensure phone index exists
DROP INDEX IF EXISTS idx_users_phone;
CREATE INDEX idx_users_phone ON users(phone);