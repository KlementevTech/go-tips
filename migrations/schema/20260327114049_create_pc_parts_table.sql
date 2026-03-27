-- +goose Up
CREATE TABLE IF NOT EXISTS pc_parts (
    id         BLOB PRIMARY KEY,
    name       TEXT NOT NULL,
    version    INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_pc_parts_created_at
    ON pc_parts (created_at DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_pc_parts_created_at;
DROP TABLE IF EXISTS pc_parts;
