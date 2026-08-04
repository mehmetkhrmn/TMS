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
    id             serial
        primary key,
    subject        varchar(255)                                               not null,
    description    text                                                       not null,
    customer_id    integer                                                    not null
        constraint fk_customer
            references customers
            on delete cascade,
    created_at     timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp with time zone default CURRENT_TIMESTAMP,
    customer_email varchar(255)                                               not null,
    status         varchar(50)              default 'open'::character varying not null
        constraint check_ticket_status
            check ((status)::text = ANY
                   ((ARRAY ['open'::character varying, 'in_progress'::character varying, 'resolved'::character varying, 'closed'::character varying])::text[])),
    constraint fk_ticket_customer_email
        foreign key (customer_id, customer_email) references customers (id, email)
            on delete cascade
);

alter table tickets
    owner to postgres;

create table answers
(
    id                serial
        primary key,
    answer            text    not null,
    representative_id integer not null
        constraint answers_agent_id_fkey
            references representatives
        constraint fk_representative
            references representatives
        constraint fk_answer_representative
            references representatives
            on delete cascade,
    ticket_id         integer not null
        references tickets
        constraint fk_ticket
            references tickets,
    answered_at       timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at        timestamp                default CURRENT_TIMESTAMP
);

alter table answers
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

