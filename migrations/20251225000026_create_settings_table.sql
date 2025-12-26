-- +goose Up
-- +goose StatementBegin
create table if not exists settings(
    user_id bigint not null,
    chat_id bigint not null,
    style text not null,
    primary key (user_id, chat_id),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists settings;
-- +goose StatementEnd
