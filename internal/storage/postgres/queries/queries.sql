-- name: CreatePcPart :one
INSERT INTO pc_parts (id, name, version, created_at, deleted_at)
VALUES (@id, @name, @version, @created_at, @deleted_at)
RETURNING *;

-- name: GetPcPart :one
SELECT id, name, version, created_at, deleted_at
FROM pc_parts
WHERE id = @id AND deleted_at IS NULL
LIMIT 1;

-- name: UpdatePcPart :one
UPDATE pc_parts
SET name = @name,
    version = @version,
    deleted_at = @deleted_at
WHERE id = @id AND version = @old_version AND deleted_at IS NULL
RETURNING *;

-- name: GetPcPartsRecent :many
SELECT id, name, version, created_at, deleted_at
FROM pc_parts
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @lim;

-- name: SoftDeletePcPart :one
UPDATE pc_parts
SET version = @version,
    deleted_at = @deleted_at
WHERE id = @id AND
    version = @old_version AND
    deleted_at IS NULL
RETURNING *;
