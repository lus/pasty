begin;

alter table "pastes" rename column "deletionToken" to "modificationToken";

commit;