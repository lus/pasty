begin;

alter table if exists "pastes" rename column "deletionToken" to "modificationToken";

commit;