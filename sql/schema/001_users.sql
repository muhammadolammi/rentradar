-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT  NOT NULL ,
    last_name TEXT  NOT NULL ,
    email TEXT UNIQUE NOT NULL ,
    phone_number TEXT UNIQUE ,
    -- role ENUM('user','agent','admin', "landlord") NOT NULL,
    role TEXT  NOT NULL ,
    password TEXT NOT NULL ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
