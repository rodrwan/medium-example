-- +goose Up
-- +goose StatementBegin
create table addresses (
    user_id       uuid not null references users(id),
    address_line  text not null,
    CHECK (address_line <> ''),
    city          text not null,
    CHECK (city <> ''),
    locality      text not null,
    CHECK (locality <> ''),
    region        text not null,
    CHECK (region <> ''),
    country       text not null,
    CHECK (country <> ''),
    postal_code   integer not null,

    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create index on addresses (user_id);

create trigger update_addresses_updated_at
before update on addresses for each row execute procedure update_updated_at_column();
-- +goose StatementEnd
-- +goose Down
drop table addresses;


