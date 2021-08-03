create table users
(
    userId            int          not null
        primary key,
    primary_group     varchar(255) not null,
    last_action       varchar(30)  null,
    last_action_value int          not null
);