-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id                  varchar   NOT NULL primary key,
    title               varchar   NOT NULL,
    description         text,
    user_id             int       NOT NULL,
    start_date          timestamp NOT NULL,
    end_date            timestamp NOT NULL,
    notification_time   int,
    day                 timestamp NOT NULL,
    week                timestamp NOT NULL,
    month               timestamp NOT NULL
);

CREATE INDEX events_day_index ON events (day);
CREATE INDEX events_week_index ON events (week);
CREATE INDEX events_month_index ON events (month);
-- +goose StatementEnd
