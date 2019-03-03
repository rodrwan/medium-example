-- +goose Up
-- +goose StatementBegin
create table users (
    id  uuid primary key default gen_random_uuid(),
    email text unique not null,
    first_name text not null,
    CHECK (first_name <> ''),
    last_name text not null,
    CHECK (last_name <> ''),
    phone text not null,
    CHECK (phone <> ''),
    birthdate timestamptz not null,

    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz
);

create trigger update_users_updated_at
before update on users for each row execute procedure update_updated_at_column();
-- +goose StatementEnd
-- +goose Down
drop table users cascade;


