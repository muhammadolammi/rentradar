-- +goose Up 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    );


-- +goose Down
DROP TABLE listings;