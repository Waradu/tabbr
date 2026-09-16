-- +goose Up
CREATE TABLE IF NOT EXISTS commands (
    command TEXT PRIMARY KEY,
    use_count INTEGER NOT NULL DEFAULT 1,
    last_used_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS exclusions (
    pattern TEXT PRIMARY KEY,
    created_at INTEGER NOT NULL
);

CREATE TEMP TABLE default_exclusions (pattern TEXT PRIMARY KEY);
INSERT INTO default_exclusions (pattern) VALUES
    ('cd *'),
    ('ls *'),
    ('*access_key*'),
    ('*api_key*'),
    ('*secret_key*'),
    ('*private_key*'),
    ('*bearer *'),
    ('*--api-key*'),
    ('*--token*'),
    ('*--password*'),
    ('*--secret*');

INSERT INTO exclusions (pattern, created_at)
SELECT pattern, unixepoch() FROM default_exclusions WHERE true
ON CONFLICT(pattern) DO NOTHING;

DROP TABLE default_exclusions;

CREATE INDEX IF NOT EXISTS commands_last_used_at_idx ON commands(last_used_at);
