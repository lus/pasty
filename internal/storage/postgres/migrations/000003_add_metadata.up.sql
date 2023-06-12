begin;

alter table "pastes" add column "metadata" jsonb;

commit;