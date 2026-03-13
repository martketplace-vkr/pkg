CREATE TABLE if not exists effect_log (
                            idempotency_key text NOT NULL,
                            effect text NOT NULL,
                            payload jsonb,
                            created_at timestamptz NOT NULL DEFAULT now(),
                            PRIMARY KEY (idempotency_key, effect)
);

CREATE INDEX IF NOT EXISTS idx_inbox_ready_null_lock
    ON public.inbox (created_at, id)
    WHERE processed_at IS NULL
    AND state = 'new'
    AND attempts < 1
    AND locked_until IS NULL;

CREATE INDEX IF NOT EXISTS idx_inbox_processed_at_id
    ON public.inbox (processed_at, id);

CREATE UNIQUE INDEX IF NOT EXISTS effect_log_idempotency_effect_uidx
    ON public.effect_log (idempotency_key, effect);

CREATE INDEX IF NOT EXISTS effect_log_created_at_idx
    ON public.effect_log (created_at);

CREATE INDEX IF NOT EXISTS idx_inbox_archive_deleted_at_id
    ON public.inbox_archive (deleted_at, id);

CREATE UNIQUE INDEX IF NOT EXISTS inbox_idempotency_key_uidx
    ON public.inbox (idempotency_key);

