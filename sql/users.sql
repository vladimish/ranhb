-- auto-generated definition
create table users
(
    userId                   bigint               not null
        primary key,
    primary_group            varchar(255)         not null,
    last_action              varchar(30)          null,
    last_action_value        int                  not null,
    is_privacy_accepted      tinyint(1) default 0 not null,
    is_terms_of_use_accepted tinyint(1) default 0 not null
);