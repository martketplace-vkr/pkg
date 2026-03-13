ALTER TABLE outbox
DROP CONSTRAINT IF EXISTS outbox_idempotency_key_unique;
DROP INDEX IF EXISTS outbox_idempotency_key_uidx;
