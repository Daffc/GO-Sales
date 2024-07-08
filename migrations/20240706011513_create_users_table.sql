-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    CONSTRAINT UC_Email UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
