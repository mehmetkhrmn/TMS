CREATE TYPE activity_action AS ENUM (
    'ticket_created',
    'ticket_updated',
    'status_changed',
    'answer_created',
    'answer_updated',
    'assignment_revoked',
    'assignment_granted',
    'message_replied'
    'message_created'
    'ticket_granted'
    'ticket_revoked'
    );

CREATE TABLE representatives
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email      VARCHAR(255) NOT NULL UNIQUE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_customers_id_email UNIQUE (id, email)
);

CREATE TABLE tickets
(
    id          SERIAL PRIMARY KEY,
    subject     VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    customer_id INTEGER NOT NULL
        CONSTRAINT fk_customer
            REFERENCES customers
            ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    status VARCHAR(50) DEFAULT 'open' NOT NULL
        CONSTRAINT check_ticket_status
            CHECK (
                status IN (
                           'open',
                           'in_progress',
                           'resolved',
                           'closed'
                    )
                ),

    category VARCHAR(50) DEFAULT 'other' NOT NULL
        CONSTRAINT check_ticket_category
            CHECK (
                category IN (
                             'technical',
                             'billing',
                             'account',
                             'bug',
                             'other'
                    )
                )
);

CREATE TABLE auth_users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    role VARCHAR(20) NOT NULL
        CONSTRAINT auth_users_role_check
            CHECK (
                role IN (
                         'admin',
                         'representative',
                         'customer'
                    )
                ),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    entity_id    INTEGER
);

CREATE TABLE activity_logs
(
    id         SERIAL PRIMARY KEY,

    ticket_id  INTEGER NOT NULL
        REFERENCES tickets
            ON DELETE CASCADE,

    user_id    INTEGER
        CONSTRAINT fk_activity_user
            REFERENCES auth_users,

    action     activity_action NOT NULL,

    field_name VARCHAR(50),
    old_value  TEXT,
    new_value  TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE ticket_assignments
(
    id                SERIAL PRIMARY KEY,

    ticket_id         INTEGER NOT NULL
        REFERENCES tickets
            ON DELETE CASCADE,

    representative_id INTEGER NOT NULL
        REFERENCES representatives
            ON DELETE CASCADE,

    assigned_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (ticket_id, representative_id)
);

CREATE TABLE ticket_messages
(
    id         SERIAL PRIMARY KEY,

    ticket_id  INTEGER NOT NULL
        REFERENCES tickets
            ON DELETE CASCADE,

    user_id    INTEGER NOT NULL
        CONSTRAINT fk_ticket_messages_user
            REFERENCES auth_users
            ON DELETE CASCADE,

    message    TEXT NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);