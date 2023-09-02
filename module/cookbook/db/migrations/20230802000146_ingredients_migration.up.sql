BEGIN;

CREATE TABLE IF NOT EXISTS ingredients (
    id          serial      PRIMARY KEY,
    name        varchar(64) NOT NULL,
    created_at  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  varchar(64) NOT NULL,
    updated_at  timestamp   NULL,
    updated_by  varchar(64) NULL,
    is_deleted  boolean     NOT NULL DEFAULT FALSE
);

INSERT INTO ingredients (name, created_by)
VALUES
    ('Nasi', 'Naufal'),
    ('Telor', 'Naufal'),
    ('Minyak', 'Naufal'),
    ('Bawang merah', 'Naufal'),
    ('Bawang putih', 'Naufal'),
    ('Ayam', 'Naufal'),
    ('Sayur cesim', 'Naufal'),
    ('Kecap', 'Naufal'),
    ('Garam', 'Naufal');

CREATE INDEX idx_ingredients_is_deleted ON ingredients(is_deleted);

COMMIT;
