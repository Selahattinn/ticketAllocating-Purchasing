CREATE DATABASE IF NOT EXISTS testdb;
GRANT ALL ON testdb.* TO 'admin'@'%';

use testdb;

create table if not exists purchase (
    id bigint auto_increment primary key,
    user_id text not null,
    quantity int  not null,
    constraint user_id_uindex
    unique (id)
);

create table if not exists ticket (
    id bigint auto_increment primary key,
    name text not null,
    description text not null,
    quantity int  not null,
    constraint ticket_id_uindex
    unique (id)
);