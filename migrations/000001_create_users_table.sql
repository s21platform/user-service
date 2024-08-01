-- +goose Up

CREATE TABLE IF NOT EXISTS users (
                                     id serial primary key,
                                     login varchar(100) not null,
                                     uuid uuid not null,
                                     email varchar(100) not null,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS data (
                                    id bigint PRIMARY KEY,
                                    user_id integer not null,
                                    name varchar(100),
                                    surname varchar(100),
                                    birthdate date,
                                    telegram varchar(100),
                                    git varchar(100),
                                    city integer,
                                    os integer,
                                    work integer,
                                    university integer
);

-- +goose Down

DROP TABLE if exists users;
DROP TABLE if exists data;
