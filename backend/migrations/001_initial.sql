CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS articles (
    id             BIGSERIAL    PRIMARY KEY,
    wikipedia_id   INTEGER      NOT NULL UNIQUE,
    title          TEXT         NOT NULL,
    slug           TEXT         NOT NULL UNIQUE,
    content        TEXT         NOT NULL DEFAULT '',
    content_length INTEGER      NOT NULL DEFAULT 0,
    rarity_tier    VARCHAR(20)  NOT NULL DEFAULT 'common',
    summary        TEXT         NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_articles_rarity ON articles(rarity_tier);

CREATE TABLE IF NOT EXISTS quiz_questions (
    id             BIGSERIAL    PRIMARY KEY,
    article_id     BIGINT       NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    question_index SMALLINT     NOT NULL,
    question_type  VARCHAR(20)  NOT NULL,
    question_text  TEXT         NOT NULL,
    options        JSONB        NOT NULL,
    correct_answer TEXT         NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(article_id, question_index)
);

CREATE INDEX IF NOT EXISTS idx_quiz_questions_article ON quiz_questions(article_id);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id           BIGSERIAL    PRIMARY KEY,
    user_id      BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    article_id   BIGINT       NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    status       VARCHAR(20)  NOT NULL,
    score        SMALLINT     NOT NULL,
    attempted_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, article_id)
);

CREATE INDEX IF NOT EXISTS idx_attempts_user    ON quiz_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_attempts_article ON quiz_attempts(article_id);

CREATE TABLE IF NOT EXISTS ownership (
    id         BIGSERIAL    PRIMARY KEY,
    article_id BIGINT       NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    claimed_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(article_id)
);

CREATE INDEX IF NOT EXISTS idx_ownership_user ON ownership(user_id);
