create table if not exists outbox
(
 -- id - auto_increment value for the id
    id             bigserial primary key,
-- key - kafka event key
    "key"          text,
-- topic is the name of the topic that the event will be published to
    topic          text      not null,
-- payload is the event data
    payload        jsonb     not null,
-- trace_carrier is the trace context that will be used to trace the event
    trace_carrier  jsonb     not null,
-- idempotent_key is a key for identification and control of repeated events
    idempotency_key text      not null,
-- created_at is the time when the event was created
    created_at     timestamp not null default now(),
-- offset is the offset of the event in the topic
    "offset"       integer,
-- partition is the partition of the event in the topic
    "partition"    integer, 
-- sent_at is the time when the event was sent
    sent_at        timestamp,
-- locked_until is the time when the event will be unlocked
    locked_until   timestamp
);


create index outbox_sending_at_idx on outbox (sent_at);

ALTER TABLE public.outbox
    ADD COLUMN meta text DEFAULT NULL,
    ADD COLUMN event_type text NOT NULL DEFAULT '';

ALTER TABLE public.outbox ADD COLUMN last_error text DEFAULT NULL;
ALTER TABLE public.outbox ADD COLUMN attempts int DEFAULT 0;

create index idx_outbox_for_send_created_at_attempts
    on outbox (created_at, id)
    where (sent_at IS NULL) AND (attempts < 1);