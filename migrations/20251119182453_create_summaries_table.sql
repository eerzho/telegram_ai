-- +goose Up
-- +goose StatementBegin
create table if not exists summaries(
    owner_id varchar(255) not null,
    peer_id varchar(255) not null,
    text text not null,
    primary key (owner_id, peer_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists summaries;
-- +goose StatementEnd
