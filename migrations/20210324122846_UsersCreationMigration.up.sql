CREATE TABLE usersauto (
    id bigserial not null primary key,
    username varchar not null unique,
    password varchar not null
);

CREATE TABLE automobiles (
    id bigserial not null primary key,
    mark varchar not null unique,
    maxspeed integer not null,
    distance integer not null,
    handler varchar not null,
    stock varchar not null
);