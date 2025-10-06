-- +goose Up 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     user_id UUID NOT NULL,
    listing_id UUID NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- ENUM('pending','sent','failed')
    status TEXT NOT NULL,

     CONSTRAINT fk_notifications_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
     CONSTRAINT fk_notifications_listing
        FOREIGN KEY (listing_id)
        REFERENCES listings(id)
        ON DELETE CASCADE
    
    );


-- +goose Down
DROP TABLE notifications;