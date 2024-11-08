-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS my_user(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(1024) NOT NULL,
    birthday DATE 
);

ALTER TABLE my_user
ALTER COLUMN birthday DROP NOT NULL; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE my_user;
-- +goose StatementEnd
