CREATE INDEX  idx_inbox_head_key_created_id
    ON inbox (key, created_at, id)
    WHERE processed_at IS NULL
      AND attempts < 1;
