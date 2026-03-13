CREATE INDEX  idx_outbox_head_key_created_id
    ON outbox (key, created_at, id)
    WHERE sent_at IS NULL
      AND attempts < 1;
