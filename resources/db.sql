-- we don't know how to generate root <with-no-name> (class Root) :(
grant alter, alter routine, create, create routine, create tablespace, create temporary tables, create user, create view, delete, drop, event, execute, file, index, insert, lock tables, process, references, reload, replication client, replication slave, select, show databases, show view, shutdown, super, trigger, update, grant option on *.* to root;

create table users_db.users
(
    id bigint auto_increment
        primary key,
    first_name varchar(45) null,
    last_name varchar(45) null,
    email varchar(45) not null,
    date_created datetime not null,
    status varchar(45) default 'active' not null,
    password varchar(256) not null,
    constraint email_unique
        unique (email)
);



