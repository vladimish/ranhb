create table payments(
    id int auto_increment,
    payment_id VARCHAR(36) NOT NULL,
    payment_time int,
    user bigint(20) NOT NULL,
    constraint premium_id_uindex
        unique (id),
    constraint `user_id`
        foreign key (user) references users(userId)
        on delete restrict
        on update restrict
);
