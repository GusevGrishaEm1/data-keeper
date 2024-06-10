-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table if not exists "data" (
    uuid uuid primary key,
    content bytea not null,
    content_type varchar(255) not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

create index if not exists data_idx on "data" (created_by);

create table if not exists "user_files" (
    uuid uuid primary key,
    content bytea not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "user_files";
DROP INDEX IF EXISTS data_idx;
DROP TABLE IF EXISTS "data";