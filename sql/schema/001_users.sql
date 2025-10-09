-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT  NOT NULL ,
    last_name TEXT  NOT NULL ,
    email TEXT UNIQUE NOT NULL ,
    phone_number TEXT UNIQUE , 
    role TEXT  NOT NULL ,
    password TEXT NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    company_name TEXT UNIQUE ,
    verified BOOLEAN NOT NULL DEFAULT false,
    rating FLOAT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE users;
