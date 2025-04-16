-- +goose Up
-- +goose StatementBegin
CREATE TYPE role AS ENUM ('employee', 'moderator');
CREATE TYPE city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
CREATE TYPE status AS ENUM ('in_progress', 'close');
CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS product_type;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS city;
DROP TYPE IF EXISTS role;
-- +goose StatementEnd
