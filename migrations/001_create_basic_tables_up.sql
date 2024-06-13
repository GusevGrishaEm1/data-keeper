-- +goose Up
create table if not exists "data" (
    uuid uuid primary key,
    content bytea not null,
    content_type varchar(255) not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

create index if not exists data_idx on "data" (created_by);

create table if not exists "user_file" (
    uuid uuid primary key,
    content bytea not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

-- +goose Down
DROP TABLE IF EXISTS "user_file";
DROP INDEX IF EXISTS data_idx;
DROP TABLE IF EXISTS "data";