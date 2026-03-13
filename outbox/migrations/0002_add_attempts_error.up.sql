alter table outbox
    add column if not exists attempts int default 0
    add column if not exists last_error text default null;