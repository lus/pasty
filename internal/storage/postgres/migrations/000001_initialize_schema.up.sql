begin;

create table "pastes" (
    "id" text not null,
    "content" text not null,
    "deletionToken" text not null,
    "created" bigint not null,
    "autoDelete" boolean not null,
    primary key ("id")
);

commit;