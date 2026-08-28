CREATE TABLE IF NOT EXISTS spots (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id             VARCHAR(255) UNIQUE NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    address              TEXT NOT NULL,
    latitude             DECIMAL(10, 7) NOT NULL,
    longitude            DECIMAL(10, 7) NOT NULL,
    first_plan_id        UUID NOT NULL REFERENCES plans(id),
    first_plan_is_public BOOLEAN NOT NULL DEFAULT FALSE
);
