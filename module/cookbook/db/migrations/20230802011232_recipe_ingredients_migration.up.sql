BEGIN;

CREATE TABLE IF NOT EXISTS recipe_ingredients (
    id                      bigserial       PRIMARY KEY,
    recipe_id               int             NOT NULL REFERENCES recipes,
    ingredient_id           int             NOT NULL REFERENCES ingredients,
    ingredient_name         varchar(64)     NOT NULL,
    ingredient_unit_name    varchar(64)     NOT NULL,
    amount                  decimal         NULL,
    ordering_index          int             NOT NULL,
    notes                   varchar(255)    NULL,
    created_at              timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              varchar(64)     NOT NULL,
    updated_at              timestamp       NULL,
    updated_by              varchar(64)     NULL,
    is_deleted              boolean         NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_recipes_recipe_id_ingredient_id_is_deleted ON recipe_ingredients(recipe_id, ingredient_id, is_deleted);

COMMIT;