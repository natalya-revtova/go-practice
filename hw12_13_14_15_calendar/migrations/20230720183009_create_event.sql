-- +goose Up
-- +goose StatementBegin
CREATE TABLE events 
(
    id          varchar   NOT NULL,
    title       varchar   NOT NULL,
    description text,
    user_id     int       NOT NULL,
    start_date  timestamp NOT NULL,
    end_date    timestamp NOT NULL,
    notification_time int
);
-- +goose StatementEnd
