alter table public.inbox
add column if not exists event_type text,
add column if not exists meta text default null;
