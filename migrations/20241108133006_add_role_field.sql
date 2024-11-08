-- +goose Up
-- +goose StatementBegin
ALTER TABLE my_user
ADD COLUMN role VARCHAR(256) NOT NULL DEFAULT 'client'; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
