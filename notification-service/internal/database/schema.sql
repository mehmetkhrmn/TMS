create table notifications
(
    id                serial
        primary key,
    ticket_id         integer                             not null,
    recipient_user_id integer                             not null,
    type              varchar(50)                         not null,
    message           text                                not null,
    created_at        timestamp default CURRENT_TIMESTAMP not null,
    actor_user_id     integer,
    occurred_at       timestamp with time zone,
    event_id          uuid
        unique
);


