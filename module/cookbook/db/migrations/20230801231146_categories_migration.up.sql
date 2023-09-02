BEGIN;

CREATE TABLE IF NOT EXISTS categories (
    id          serial      PRIMARY KEY,
    name        varchar(64) NOT NULL,
    created_at  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  varchar(64) NOT NULL,
    updated_at  timestamp   NULL,
    updated_by  varchar(64) NULL,
    is_deleted  boolean     NOT NULL DEFAULT FALSE
);

INSERT INTO categories (name, created_by)
VALUES
    ('Main Course', 'Naufal'),
    ('Beverages', 'Naufal'),
    ('Dessert', 'Naufal');

CREATE INDEX idx_categories_is_deleted ON categories(is_deleted);

COMMIT;
