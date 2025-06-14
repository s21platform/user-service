-- +goose Up

ALTER TABLE users
    ADD CONSTRAINT unique_uuid UNIQUE (uuid);

ALTER TABLE data
    ADD COLUMN temp_uuid UUID;

ALTER TABLE skills
    ADD COLUMN temp_uuid UUID;

ALTER TABLE hobbies
    ADD COLUMN temp_uuid UUID;

ALTER TABLE posts
    ADD COLUMN temp_uuid UUID;

UPDATE data
SET temp_uuid = (SELECT uuid FROM users WHERE users.id = data.user_id);

UPDATE skills
SET temp_uuid = (SELECT uuid FROM users WHERE users.id = skills.user_id);

UPDATE hobbies
SET temp_uuid = (SELECT uuid FROM users WHERE users.id = hobbies.user_id);

UPDATE posts
SET temp_uuid = (SELECT uuid FROM users WHERE users.uuid = posts.user_id);

ALTER TABLE data
DROP CONSTRAINT fk_data_user_id;
ALTER TABLE data
ALTER COLUMN user_id TYPE UUID USING temp_uuid;
ALTER TABLE data
    RENAME COLUMN user_id TO user_uuid;
ALTER TABLE data
DROP COLUMN temp_uuid;

ALTER TABLE skills
DROP CONSTRAINT fk_skill_user_id;
ALTER TABLE skills
ALTER COLUMN user_id TYPE UUID USING temp_uuid;
ALTER TABLE skills
    RENAME COLUMN user_id TO user_uuid;
ALTER TABLE skills
DROP COLUMN temp_uuid;

ALTER TABLE hobbies
DROP CONSTRAINT fk_hobby_user_id;
ALTER TABLE hobbies
ALTER COLUMN user_id TYPE UUID USING temp_uuid;
ALTER TABLE hobbies
    RENAME COLUMN user_id TO user_uuid;
ALTER TABLE hobbies
DROP COLUMN temp_uuid;

ALTER TABLE posts
ALTER COLUMN user_id TYPE UUID USING temp_uuid;
ALTER TABLE posts
    RENAME COLUMN user_id TO user_uuid;
ALTER TABLE posts
DROP COLUMN temp_uuid;

ALTER TABLE data
    ADD CONSTRAINT fk_data_user_uuid FOREIGN KEY (user_uuid) REFERENCES users(uuid);

ALTER TABLE skills
    ADD CONSTRAINT fk_skill_user_uuid FOREIGN KEY (user_uuid) REFERENCES users(uuid);

ALTER TABLE hobbies
    ADD CONSTRAINT fk_hobby_user_uuid FOREIGN KEY (user_uuid) REFERENCES users(uuid);

ALTER TABLE posts
    ADD CONSTRAINT fk_post_user_uuid FOREIGN KEY (user_uuid) REFERENCES users(uuid);
