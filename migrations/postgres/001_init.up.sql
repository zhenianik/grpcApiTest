create table users
(
    id    serial
        constraint users_pkey primary key,
    name  text not null UNIQUE,
    email text not null UNIQUE
);