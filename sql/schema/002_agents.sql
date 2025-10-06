-- +goose Up 
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    company_name TEXT UNIQUE,
    verified BOOLEAN DEFAULT false,
    rating FLOAT DEFAULT 0,


    CONSTRAINT fk_agent_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE


);
-- +goose Down
DROP TABLE agents;
