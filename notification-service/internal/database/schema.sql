CREATE TABLE notifications (
                               id SERIAL PRIMARY KEY,
                               ticket_id INTEGER NOT NULL,
                               user_id INTEGER NOT NULL,
                               type VARCHAR(50) NOT NULL,
                               message TEXT NOT NULL,
                               created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);