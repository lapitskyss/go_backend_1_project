CREATE TABLE IF NOT EXISTS links
(
    id          bigserial PRIMARY KEY,
    url         varchar(10000) not null,
    hash        varchar(20)  not null,
    created_at  TIMESTAMP DEFAULT current_timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS links_hashx
    ON links (hash);
