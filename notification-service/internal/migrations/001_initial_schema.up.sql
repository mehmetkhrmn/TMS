CREATE TABLE notifications
(
    id                SERIAL PRIMARY KEY,
    ticket_id         INTEGER NOT NULL,
    recipient_user_id INTEGER NOT NULL,
    type              VARCHAR(50) NOT NULL,
    message           TEXT NOT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    actor_user_id     INTEGER,
    occurred_at       TIMESTAMP WITH TIME ZONE,
    event_id          UUID UNIQUE
);