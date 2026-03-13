create table if not exists outbox_sending(
    idempotency_key text not null,
    topic text not null,
    event_type text not null,
    payload jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (idempotency_key, topic, event_type)
);

CREATE INDEX IF NOT EXISTS idx_outbox_sending_created_at ON outbox_sending(created_at);

