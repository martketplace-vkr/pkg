CREATE UNIQUE INDEX IF NOT EXISTS inbox_idempotency_key_uidx
    ON inbox (idempotency_key);

ALTER TABLE inbox
    ADD CONSTRAINT inbox_idempotency_key_unique
    UNIQUE USING INDEX inbox_idempotency_key_uidx;
