-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank
(
    id varchar(255),
    bank_code varchar(50) not null,
    bank_name varchar(255) not null,
    bank_number varchar(50) not null,
    content varchar(500),
    is_default bool not null default false,
    customer_name varchar(100) not null,
    owner varchar(255) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    primary key (id)
);

CREATE OR REPLACE FUNCTION update_bank_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_bank_updated_at BEFORE UPDATE
    ON bank FOR EACH ROW EXECUTE PROCEDURE
    update_bank_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table bank;;
-- +goose StatementEnd
