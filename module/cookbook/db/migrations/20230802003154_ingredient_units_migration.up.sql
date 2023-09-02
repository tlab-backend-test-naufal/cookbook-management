BEGIN;

CREATE TABLE IF NOT EXISTS ingredient_units (
    id          serial      PRIMARY KEY,
    name        varchar(64) NOT NULL,
    created_at  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  varchar(64) NOT NULL,
    updated_at  timestamp   NULL,
    updated_by  varchar(64) NULL,
    is_deleted  boolean     NOT NULL DEFAULT FALSE
);

INSERT INTO ingredient_units (name, created_by)
VALUES
    ('piring', 'Naufal'),
    ('buah', 'Naufal'),
    ('siung', 'Naufal'),
    ('helai', 'Naufal');

CREATE INDEX idx_ingredient_units_is_deleted ON ingredient_units(is_deleted);

COMMIT;
