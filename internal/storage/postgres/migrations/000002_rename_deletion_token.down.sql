begin;

alter table "pastes" rename column "modificationToken" to "deletionToken";

commit;