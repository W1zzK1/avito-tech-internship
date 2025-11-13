-- Создаем таблицу команд
CREATE TABLE IF NOT EXISTS teams
(
    id   TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

-- Создаем таблицу пользователей (связь многие-к-одному с teams)
CREATE TABLE IF NOT EXISTS users
(
    id        TEXT PRIMARY KEY,
    username  TEXT    NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_id   TEXT    NOT NULL REFERENCES teams (id) ON DELETE CASCADE
);

-- Создаем таблицу pull requests
CREATE TABLE IF NOT EXISTS pull_requests
(
    id         TEXT PRIMARY KEY,
    name       TEXT      NOT NULL,
    author_id  TEXT      NOT NULL REFERENCES users (id),
    status     TEXT      NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at  TIMESTAMP
);

-- Создаем таблицу для связи многие-ко-многим между PR и ревьюерами
CREATE TABLE IF NOT EXISTS pull_request_reviewers
(
    pull_request_id TEXT NOT NULL REFERENCES pull_requests (id) ON DELETE CASCADE,
    user_id         TEXT NOT NULL REFERENCES users (id),
    PRIMARY KEY (pull_request_id, user_id)
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users (team_id, is_active);
CREATE INDEX IF NOT EXISTS idx_prs_author_status ON pull_requests (author_id, status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pull_request_reviewers (user_id);