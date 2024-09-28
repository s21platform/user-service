-- +goose Up

CREATE TABLE IF NOT EXISTS users
(
    id               serial primary key,
    login            varchar(100) not null,
    uuid             uuid         not null,
    email            varchar(100) not null,
    last_avatar_link TEXT         not null,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS data
(
    id            serial PRIMARY KEY,
    user_id       integer not null,
    name          varchar(100),
    surname       varchar(100),
    birthdate     date,
    phone         varchar(15),
    telegram      varchar(100),
    git           varchar(100),
    city_id       integer,
    os_id         integer,
    work_id       integer,
    university_id integer,
    CONSTRAINT fk_data_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS skills
(
    id serial PRIMARY KEY,
    user_id bigint not null,
    skill_id bigint not null,
    CONSTRAINT fk_skill_user_id FOREIGN KEY (user_id) REFERENCES users(id)

);

CREATE TABLE IF NOT EXISTS hobbies
(
    id serial PRIMARY KEY,
    user_id bigint not null,
    hobby_id bigint not null,
    CONSTRAINT fk_hobby_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE if exists data;
DROP TABLE if exists skills;
DROP TABLE if exists hobbies;
DROP TABLE if exists users;
