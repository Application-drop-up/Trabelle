CREATE TABLE IF NOT EXISTS notes (
    id          SERIAL PRIMARY KEY,
    pin_id      INTEGER NOT NULL REFERENCES pins(id),
    content     TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
