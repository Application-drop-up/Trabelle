CREATE TABLE IF NOT EXISTS plans (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    share_token  VARCHAR(255) UNIQUE NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
