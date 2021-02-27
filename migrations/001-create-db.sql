DROP TABLE IF EXISTS ad;

CREATE TABLE IF NOT EXISTS ad
(
    id       UUID PRIMARY KEY,
    content  TEXT,
    created  TIMESTAMP NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at   TIMESTAMP NOT NULL
);

