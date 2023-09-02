BEGIN;

CREATE TABLE IF NOT EXISTS recipes (
    id          bigserial       PRIMARY KEY,
    name        varchar(64)     NOT NULL,
    description varchar(255)    NULL,
    category_id int             NOT NULL REFERENCES categories,
    created_at  timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  varchar(64)     NOT NULL,
    updated_at  timestamp       NULL,
    updated_by  varchar(64)     NULL,
    is_deleted  boolean         NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_recipes_category_id_is_deleted ON recipes(category_id, is_deleted);

COMMIT;