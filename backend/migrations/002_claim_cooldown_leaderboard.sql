CREATE TABLE IF NOT EXISTS claim_cooldowns (
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_claim TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id)
);
