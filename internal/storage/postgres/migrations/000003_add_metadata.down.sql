begin;

alter table if exists "pastes" drop column "metadata";

commit;