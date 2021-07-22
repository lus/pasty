begin;

alter table if exists "pastes" add column "metadata" jsonb not null;

commit;