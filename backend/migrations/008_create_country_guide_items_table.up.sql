CREATE TABLE IF NOT EXISTS country_guide_items (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_guide_id UUID NOT NULL REFERENCES country_guides(id) ON DELETE CASCADE,
    category         VARCHAR(50) NOT NULL,
    title            VARCHAR(255) NOT NULL,
    description      TEXT,
    url              VARCHAR(500),
    is_mandatory     BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order       INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_country_guide_items_country_guide_id ON country_guide_items(country_guide_id);
