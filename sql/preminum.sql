create table premium
(
    id         int auto_increment,
    user_id    bigint null,
    begin_date int    null,
    end_date   int    null,
    constraint premium_id_uindex
        unique (id),
    constraint premium_ibfk_1
        foreign key (user_id) references ranh.users (userId)
);

create index user_id
    on premium (user_id);
