-- +goose Up
-- +goose StatementBegin
create table if not exists summaries(
    chat_id varchar(255) primary key,
    text text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exitst summaries;
-- +goose StatementEnd
