CREATE TYPE state AS ENUM ('new', 'processing', 'done', 'error');

ALTER TABLE public.inbox
    ADD COLUMN state state DEFAULT 'done' NOT NULL;

create index idx_inbox_for_processing_created_at_attempts
    on public.inbox (created_at, id)
    where (processed_at IS NULL) AND (attempts < 1) AND (state='new');