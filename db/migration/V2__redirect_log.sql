CREATE TABLE IF NOT EXISTS redirects_logs
(
    id          bigserial PRIMARY KEY,
    link_id     bigint NOT NULL REFERENCES links(id),
    user_agent  varchar(1000) DEFAULT NULL,
    created_at  TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS redirects_logs_link_idx
    ON redirects_logs (link_id)
