CREATE TABLE IF NOT EXISTS spots (
    id               SERIAL PRIMARY KEY,
    google_place_id  VARCHAR(255) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    address          TEXT,
    latitude         DECIMAL(10, 7) NOT NULL,
    longitude        DECIMAL(10, 7) NOT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW()
);
