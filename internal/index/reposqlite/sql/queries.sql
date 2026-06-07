-- name: InsertSymbol :exec
INSERT INTO symbols (id, name, kind, package, file, signature)
VALUES (?, ?, ?, ?, ?, ?);

-- name: InsertRelation :exec
INSERT INTO relations (from_id, to_id, edge_kind)
VALUES (?, ?, ?);

-- name: InsertMethod :exec
INSERT INTO methods (parent_id, name, signature, kind)
VALUES (?, ?, ?, ?);

-- name: InsertFile :exec
INSERT INTO files (path, lang, hash, summary)
VALUES (?, ?, ?, ?);

-- name: DeleteAllSymbols :exec
DELETE FROM symbols;

-- name: DeleteAllRelations :exec
DELETE FROM relations;

-- name: DeleteAllMethods :exec
DELETE FROM methods;

-- name: DeleteAllFiles :exec
DELETE FROM files;

-- name: ListSymbols :many
SELECT id, name, kind, package, file, signature
FROM symbols
ORDER BY name;

-- name: GetCallsFrom :many
SELECT s.name
FROM relations r
JOIN symbols s ON s.id = r.to_id
WHERE r.from_id = ? AND r.edge_kind = 'call'
ORDER BY s.name;

-- name: GetCallersOf :many
SELECT s.name
FROM relations r
JOIN symbols s ON s.id = r.from_id
WHERE r.to_id = ? AND r.edge_kind = 'call'
ORDER BY s.name;

-- name: GetImplementationsOf :many
SELECT s.name
FROM relations r
JOIN symbols s ON s.id = r.to_id
WHERE r.from_id = ? AND r.edge_kind = 'implement'
ORDER BY s.name;

-- name: GetMethodsOf :many
SELECT name
FROM methods
WHERE parent_id = ?
ORDER BY name;
