-- +goose Up 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    agent_id UUID NOT NULL,
    title TEXT  NOT NULL ,
    description TEXT NOT NULL,
    -- Annual or Monthly
    rent_type TEXT NOT NULL, 
    price BIGINT NOT NULL,
    -- City/Area
    location TEXT NOT NULL,
    latitude FLOAT ,
    longtitude FLOAT,
    -- shared/selcontained/one bedroom flat etc.
    house_type TEXT NOT NULL,
    -- if the apartment is verified by admin
    verified BOOLEAN NOT NULL  DEFAULT false,
    --  images JSONB NOT NULL,
    images JSON NOT NULL,
    --  ENUM('active','inactive','rented')
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_listings_agent
        FOREIGN KEY (agent_id)
        REFERENCES agents(id)
        ON DELETE CASCADE
    
    );


-- +goose Down
DROP TABLE listings;