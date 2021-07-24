create table users
(
    userId int not null,
    primaryGroup varchar(255) null,
    constraint users_pk
        primary key (userId)
);
