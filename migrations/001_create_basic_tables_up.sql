-- +goose Up
create table if not exists user_data (
    uuid uuid primary key,
    content bytea not null,
    content_type varchar(255) not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

create index if not exists data_idx on user_data (created_by);

create table if not exists file_repository (
    uuid uuid primary key,
    content bytea not null,
    created_at timestamp not null,
    created_by varchar(255) not null
);

-- +goose Down
DROP TABLE IF EXISTS file_repository;
DROP INDEX IF EXISTS data_idx;
DROP TABLE IF EXISTS user_data;