CREATE UNIQUE INDEX IF NOT EXISTS outbox_idempotency_key_uidx
    ON outbox (idempotency_key);

ALTER TABLE outbox
    ADD CONSTRAINT outbox_idempotency_key_unique
    UNIQUE USING INDEX outbox_idempotency_key_uidx;
