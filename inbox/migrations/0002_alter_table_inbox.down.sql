alter table inbox 
drop column if exists meta,
drop column if exists event_type;