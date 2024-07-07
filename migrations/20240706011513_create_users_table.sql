-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name text,
    email text,
    password text,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
