begin;

alter table if exists "pastes" drop column "autoDelete";

commit;