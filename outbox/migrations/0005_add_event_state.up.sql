CREATE TYPE state AS ENUM ('new', 'processing', 'done', 'error');

ALTER TABLE outbox
    ADD COLUMN state state DEFAULT 'done' NOT NULL;