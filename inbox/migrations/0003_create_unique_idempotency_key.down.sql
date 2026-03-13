ALTER TABLE inbox
DROP CONSTRAINT IF EXISTS inbox_idempotency_key_unique;
DROP INDEX IF EXISTS inbox_idempotency_key_uidx;
