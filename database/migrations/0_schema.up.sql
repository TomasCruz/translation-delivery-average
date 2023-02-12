CREATE TABLE IF NOT EXISTS events (
    event_id              CHAR(20) NOT NULL PRIMARY KEY,
    event_name            VARCHAR(64) NOT NULL,       
    event_ts              TIMESTAMP NOT NULL,
    payload               VARCHAR(512) NOT NULL
);

CREATE INDEX events_event_ts_idx ON events (event_ts);
CREATE INDEX events_event_name_idx ON events (event_name);
