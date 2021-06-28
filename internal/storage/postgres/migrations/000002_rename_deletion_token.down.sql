begin;

alter table if exists "pastes" rename column "modificationToken" to "deletionToken";

commit;