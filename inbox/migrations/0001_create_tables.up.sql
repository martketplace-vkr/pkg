create table if not exists inbox
(
-- id - auto_increment value for the id
    id             bigserial primary key,
-- key - kafka event key
    "key"          text,
-- entity_id is a identifier for the entity that the event is related to (order_id for example)
-- to group events by entity
    topic          text      not null,
-- payload is the event data
    payload        jsonb     not null,
-- trace_carrier is the trace context that will be used to trace the event
    trace_carrier  jsonb     not null,
-- idempotent_key is a key for identification and control of repeated events
    idempotency_key text      not null,
-- created_at is the time when the event was created
    created_at     timestamptz not null default now(),
-- processed_at is the time when the event was sent
    processed_at        timestamptz,
-- locked_until is the time when the event will be unlocked
    locked_until   timestamptz,
-- error is the process error
    error text,
-- attempts amount of process attempts
    attempts int default 0
);

-- index to search by processed_at
create index inbox_processed_at_idx on public.inbox (processed_at);
