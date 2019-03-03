-- +goose Up
-- +goose StatementBegin
create function update_updated_at_column()
returns trigger as $$
  begin
    new.updated_at = now();
    return new;
  end;
$$ language plpgsql;

create extension citext;
create extension pgcrypto;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function update_updated_at_column();
drop extension citext;
drop extension pgcrypto;
-- +goose StatementEnd
