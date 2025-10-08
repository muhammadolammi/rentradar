-- +goose Up 

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    min_price BIGINT NOT NULL,
    max_price BIGINT NOT NULL,
    location TEXT NOT NULL,
    property_type_id UUID NOT NULL, 
    --  ENUM('email','whatsapp','sms')
    contact_method TEXT NOT NULL,

    CONSTRAINT fk_alerts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE ,
    CONSTRAINT fk_alerts_property_type
        FOREIGN KEY (property_type_id)
        REFERENCES property_types(id)
        ON DELETE CASCADE
    

    
    );


-- +goose Down
DROP TABLE alerts;