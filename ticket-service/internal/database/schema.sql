create table representatives
(
    id         serial
        primary key,
    name       varchar(255) not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);

alter table representatives
    owner to postgres;

create table customers
(
    id         serial
        primary key,
    name       varchar(255) not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    email      varchar(255) not null
        unique,
    updated_at timestamp default CURRENT_TIMESTAMP,
    constraint uk_customers_id_email
        unique (id, email)
);

alter table customers
    owner to postgres;

create table tickets
(
    id          serial
        primary key,
    subject     varchar(255)                                                not null,
    description text                                                        not null,
    customer_id integer                                                     not null
        constraint fk_customer
            references customers
            on delete cascade,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    status      varchar(50)              default 'open'::character varying  not null
        constraint check_ticket_status
            check ((status)::text = ANY
                   ((ARRAY ['open'::character varying, 'in_progress'::character varying, 'resolved'::character varying, 'closed'::character varying])::text[])),
    category    varchar(50)              default 'other'::character varying not null
        constraint check_ticket_category
            check ((category)::text = ANY
                   ((ARRAY ['technical'::character varying, 'billing'::character varying, 'account'::character varying, 'bug'::character varying, 'other'::character varying])::text[]))
);

alter table tickets
    owner to postgres;

create table auth_users
(
    id            serial
        primary key,
    username      varchar(50)                         not null
        unique,
    password_hash text                                not null,
    role          varchar(20)                         not null
        constraint auth_users_role_check
            check ((role)::text = ANY
        ((ARRAY ['admin'::character varying, 'representative'::character varying, 'customer'::character varying])::text[])),
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    updated_at    timestamp default CURRENT_TIMESTAMP not null,
    entity_id     integer
);

alter table auth_users
    owner to postgres;

create table activity_logs
(
    id         serial
        primary key,
    ticket_id  integer                             not null
        references tickets
            on delete cascade,
    user_id    integer
        constraint fk_activity_user
            references auth_users,
    action     varchar(50)                         not null,
    field_name varchar(50),
    old_value  text,
    new_value  text,
    created_at timestamp default CURRENT_TIMESTAMP not null
);

alter table activity_logs
    owner to postgres;

create table ticket_assignments
(
    id                serial
        primary key,
    ticket_id         integer not null
        references tickets
            on delete cascade,
    representative_id integer not null
        references representatives
            on delete cascade,
    assigned_at       timestamp with time zone default CURRENT_TIMESTAMP,
    unique (ticket_id, representative_id)
);

alter table ticket_assignments
    owner to postgres;

create table ticket_messages
(
    id         serial
        primary key,
    ticket_id  integer not null
        references tickets
            on delete cascade,
    user_id    integer not null
        references auth_users
            on delete cascade
        constraint fk_ticket_messages_user
            references auth_users
            on delete cascade,
    message    text    not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP
);

alter table ticket_messages
    owner to postgres;

