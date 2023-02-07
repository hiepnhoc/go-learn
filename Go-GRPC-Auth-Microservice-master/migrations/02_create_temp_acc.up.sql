DROP TABLE IF EXISTS acc CASCADE;

CREATE TABLE acc
(
    user_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)              NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(32)              NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE       NOT NULL CHECK ( email <> '' ),
    password   VARCHAR(250)             NOT NULL CHECK ( octet_length(password) <> 0 ),
    role       role                     NOT NULL DEFAULT 'user',
    avatar     VARCHAR(250),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);