-- +goose Up 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- the user email or phone number
    contact TEXT NOT NULL,
    contact_method TEXT NOT NULL,
    -- ENUM('pending','sent','failed')
    status TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
     CONSTRAINT fk_notifications_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
    
    );


-- +goose Down
DROP TABLE IF EXISTS notifications;