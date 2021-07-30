begin;

alter table if exists "pastes" add column "autoDelete" boolean;

commit;