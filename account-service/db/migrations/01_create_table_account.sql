-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts
(
    id varchar,
    name varchar(255) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
                                 primary key (id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accounts;
-- +goose StatementEnd
