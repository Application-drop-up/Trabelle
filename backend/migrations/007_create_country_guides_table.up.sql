CREATE TABLE IF NOT EXISTS country_guides (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) UNIQUE NOT NULL,
    country_name VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
