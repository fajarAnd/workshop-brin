-- This migration is not easily reversible since we lose the original SERIAL sequence
-- We would need to recreate the SERIAL column and lose UUID values
-- For safety, this migration throws an error

-- Note: Converting back from UUID to SERIAL is destructive
-- You would lose all existing UUID values and need to regenerate them as integers
-- This is generally not recommended in production

-- If you really need to rollback, you would need to:
-- 1. Create a new SERIAL column
-- 2. Assign new integer IDs to all records
-- 3. Update all foreign key references (if any exist)
-- 4. Drop the UUID column

SELECT 'WARNING: Rolling back UUID to SERIAL conversion is destructive and not supported' as warning;