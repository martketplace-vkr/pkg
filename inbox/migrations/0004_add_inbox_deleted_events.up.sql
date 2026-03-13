CREATE TABLE IF NOT EXISTS inbox_archive (
                         id              bigint PRIMARY KEY,
                         key             text,
                         topic           text NOT NULL,
                         payload         jsonb NOT NULL,
                         trace_carrier   jsonb NOT NULL,
                         idempotency_key text NOT NULL,
                         created_at      timestamptz NOT NULL,
                         processed_at    timestamptz,
                         locked_until    timestamptz,
                         error           text,
                         attempts        int DEFAULT 0,
                         event_type      text,
                         meta            text,
                         deleted_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_inbox_archive_deleted_at ON public.inbox_archive(deleted_at);
