CREATE TABLE IF NOT EXISTS pins (
    id          SERIAL PRIMARY KEY,
    plan_id     INTEGER NOT NULL REFERENCES plans(id),
    spot_id     INTEGER NOT NULL REFERENCES spots(id),
    category    VARCHAR(255) NOT NULL,
    color       VARCHAR(7) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
